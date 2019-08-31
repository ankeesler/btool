// Package scanner provides a type that can collect node.Node's from an FS.
package scanner

import (
	"path/filepath"
	"strings"

	"github.com/ankeesler/btool/collector"
	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Includeser

// Includeser is a type that can return a list of #include's from a given file.
type Includeser interface {
	Includes(path string) ([]string, error)
}

// Scanner will collect nodes from an FS. It is provided a root directory and
// returns all paths prefixed with the root directory.
type Scanner struct {
	fs   afero.Fs
	root string
	i    Includeser
}

// New creates a new Scanner.
func New(
	fs afero.Fs,
	root string,
	i Includeser,
) *Scanner {
	return &Scanner{
		fs:   fs,
		root: root,
		i:    i,
	}
}

// This should be a stack so that you only get your descendents' includePaths
// and libraries!
type state struct {
	ctx *collector.Ctx

	includePaths []string
	libraries    []*node.Node
	objects      []*node.Node
	cc           bool

	depth int
}

// Collect will build up a node.Node graph given a starting node.Node. It will walk
// the dependencies of the node.Node and build up a graph, or return an error if
// it runs into trouble.
func (s *Scanner) Collect(ctx *collector.Ctx, start *node.Node) error {
	_, err := s.add(start, &state{
		ctx: ctx,
	})
	return err
}

func (s *Scanner) add(n *node.Node, state *state) (bool, error) {
	if state.ctx.NS.Find(n.Name) != nil {
		return false, nil
	}

	state.depth++
	if state.depth == 100 {
		return false, errors.New("hit depth of 100, failing")
	}

	log.Debugf("adding %s", n.Name)
	state.ctx.NS.Add(n)

	var err error
	ext := filepath.Ext(n.Name)
	switch ext {
	case "":
		err = s.onExecutable(n, state)
	//case ".a":
	//	err = s.onLibrary(n, state)
	case ".o":
		err = s.onObject(n, state)
	case ".c", ".cc":
		err = s.onSource(n, state)
	case ".h":
		err = s.onHeader(n, state)
	default:
		err = errors.New("unknown file ext: " + ext)
	}

	return true, err
}

// the path currently goes like this:
//   executable -> source -> header -> source -> header ... -> object

func (s *Scanner) onExecutable(n *node.Node, state *state) error {
	if err := s.addSource(n, state); err != nil {
		return errors.Wrap(err, "add source")
	}

	for _, oN := range state.objects {
		log.Debugf("object dependency %s -> %s", n.Name, oN.Name)
		n.Dependency(oN)
	}

	for _, lN := range state.libraries {
		log.Debugf("library dependency %s -> %s", n.Name, lN.Name)
		n.Dependency(lN)
	}

	if state.cc {
		n.Resolver = state.ctx.RF.NewLinkCC()
	} else {
		n.Resolver = state.ctx.RF.NewLinkC()
	}

	return nil
}

func (s *Scanner) onObject(n *node.Node, state *state) error {
	if state.cc {
		n.Resolver = state.ctx.RF.NewCompileCC(state.includePaths)
	} else {
		n.Resolver = state.ctx.RF.NewCompileC(state.includePaths)
	}
	state.includePaths = state.includePaths[:0] // clear

	log.Debugf("adding object %s", n.Name)
	state.objects = append(state.objects, n)

	return nil
}

func (s *Scanner) onSource(n *node.Node, state *state) error {
	ext := filepath.Ext(n.Name)
	if err := s.addHeaders(n, state); err != nil {
		return errors.Wrap(err, "add headers")
	}

	object := strings.ReplaceAll(n.Name, ext, ".o")
	objectN := node.New(object)
	log.Debugf("source/object dependency %s -> %s", objectN.Name, n.Name)
	objectN.Dependency(n)
	if added, err := s.add(objectN, state); err != nil {
		return errors.Wrap(err, "add")
	} else if !added {
		return nil
	}

	return nil
}

func (s *Scanner) onHeader(n *node.Node, state *state) error {
	if err := s.addHeaders(n, state); err != nil {
		return errors.Wrap(err, "add headers")
	}

	if err := s.addSource(n, state); err != nil {
		return errors.Wrap(err, "add source")
	}

	return nil
}

func (s *Scanner) addSource(n *node.Node, state *state) error {
	var sourceC, sourceCC string
	ext := filepath.Ext(n.Name)
	if ext == "" {
		sourceC = n.Name + ".c"
		sourceCC = n.Name + ".cc"
	} else {
		sourceC = strings.ReplaceAll(n.Name, ext, ".c")
		sourceCC = strings.ReplaceAll(n.Name, ext, ".cc")
	}
	log.Debugf("considering %s and %s", sourceC, sourceCC)

	var sourceN *node.Node
	if exists, err := afero.Exists(s.fs, sourceC); err != nil {
		return errors.Wrap(err, "exists")
	} else if exists {
		sourceN = node.New(sourceC)
	} else if exists, err = afero.Exists(s.fs, sourceCC); err != nil {
		return errors.Wrap(err, "exists")
	} else if exists {
		sourceN = node.New(sourceCC)
		state.cc = true
	}

	if sourceN == nil {
		return errors.New("unknown source for node " + n.Name)
	}

	if added, err := s.add(sourceN, state); err != nil {
		return errors.Wrap(err, "add "+sourceN.Name)
	} else if !added {
		return nil
	}

	return nil
}

func (s *Scanner) addHeaders(n *node.Node, state *state) error {
	includes, err := s.i.Includes(n.Name)
	if err != nil {
		return errors.Wrap(err, "includes")
	}

	for _, include := range includes {
		log.Debugf("searching for include %s in %+v", include, state.ctx)

		var headerN *node.Node
		rootRelInclude := filepath.Join(s.root, include)
		dirRelInclude := filepath.Join(filepath.Dir(n.Name), include)
		if exists, err := afero.Exists(s.fs, rootRelInclude); err != nil {
			return errors.Wrap(err, "exists")
		} else if exists {
			headerN = node.New(rootRelInclude)
			if !contains(state.includePaths, s.root) {
				state.includePaths = append(state.includePaths, s.root)
			}
		} else if exists, err = afero.Exists(s.fs, dirRelInclude); err != nil {
			return errors.Wrap(err, "exists")
		} else if exists {
			headerN = node.New(dirRelInclude)
		} else if headerN = s.headerForInclude(include, state); headerN != nil {
			includePath := strings.ReplaceAll(headerN.Name, include, "")
			state.includePaths = append(state.includePaths, includePath)

			log.Debugf("added include path %s", includePath)

			if libraries := state.ctx.Libraries(include); libraries != nil {
				log.Debugf("adding libraries %s for include %s", libraries, include)
				state.libraries = append(state.libraries, libraries...)
			}
		}

		if headerN == nil {
			return errors.New("unknown header for include " + include)
		}

		if _, err := s.add(headerN, state); err != nil {
			return errors.Wrap(err, "add "+headerN.Name)
		}

		log.Debugf("include dependency %s -> %s", n.Name, headerN.Name)
		n.Dependency(headerN)
	}

	return nil
}

func (s *Scanner) headerForInclude(include string, state *state) *node.Node {
	for _, includePath := range state.ctx.IncludePaths() {
		header := filepath.Join(includePath, include)
		if headerN := state.ctx.NS.Find(header); headerN != nil {
			return headerN
		}
	}
	return nil
}

func contains(ss []string, s string) bool {
	for i := range ss {
		if ss[i] == s {
			return true
		}
	}
	return false
}
