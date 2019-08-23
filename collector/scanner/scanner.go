// Package scanner provides a type that can collect node.Node's from an FS.
package scanner

import (
	"path/filepath"
	"strings"

	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . NodeStore

// NodeStore is a thing that can create and find node.Node's.
type NodeStore interface {
	Add(*node.Node)
	Find(string) *node.Node
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Includeser

// Includeser is a type that can return a list of #include's from a given file.
type Includeser interface {
	Includes(path string) ([]string, error)
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . ResolverFactory

// ResolverFactory can create node.Resolver's.
type ResolverFactory interface {
	NewCompileC(includeDirs []string) node.Resolver
	NewCompileCC(includeDirs []string) node.Resolver
	NewArchive() node.Resolver
	NewLinkC() node.Resolver
	NewLinkCC() node.Resolver
	NewSymlink() node.Resolver

	NewDownload(url, sha256 string) node.Resolver
	NewUnzip(outputDir string) node.Resolver
}

// Scanner will collect nodes from an FS. It is provided a root directory and
// returns all paths prefixed with the root directory.
type Scanner struct {
	fs   afero.Fs
	root string
	ns   NodeStore
	i    Includeser
	rf   ResolverFactory
}

// New creates a new Scanner.
func New(
	fs afero.Fs,
	root string,
	ns NodeStore,
	i Includeser,
	rf ResolverFactory,
) *Scanner {
	return &Scanner{
		fs:   fs,
		root: root,
		ns:   ns,
		i:    i,
		rf:   rf,
	}
}

// This should be a stack so that you only get your descendents' includePaths
// and libraries!
type ctx struct {
	includePaths []string
	libraries    []*node.Node
	objects      []*node.Node
	cc           bool

	depth int
}

// Scan will build up a node.Node graph given a starting node.Node. It will walk
// the dependencies of the node.Node and build up a graph, or return an error if
// it runs into trouble.
func (s *Scanner) Scan(start *node.Node) error {
	_, err := s.add(start, &ctx{
		includePaths: []string{s.root},
	})
	return err
}

func (s *Scanner) add(n *node.Node, ctx *ctx) (bool, error) {
	if s.ns.Find(n.Name) != nil {
		return false, nil
	}

	ctx.depth++
	if ctx.depth == 100 {
		return false, errors.New("hit depth of 100, failing")
	}

	log.Debugf("adding %s", n.Name)
	s.ns.Add(n)

	var err error
	ext := filepath.Ext(n.Name)
	switch ext {
	case "":
		err = s.onExecutable(n, ctx)
	//case ".a":
	//	err = s.onLibrary(n, ctx)
	case ".o":
		err = s.onObject(n, ctx)
	case ".c", ".cc":
		err = s.onSource(n, ctx)
	case ".h":
		err = s.onHeader(n, ctx)
	default:
		err = errors.New("unknown file ext: " + ext)
	}

	return true, err
}

// each of these functions should do the following:
//   - find dependencies of the provided node.Node
//   - add those dependencies as nodes with s.add()
//   - add those dependencies as dependencies with n.Dependency()
//   - set n.Resolver
//   - update ctx with any related stuff
//
// the path currently goes like this:
//   executable -> source -> header -> source -> header ... -> object

func (s *Scanner) onExecutable(n *node.Node, ctx *ctx) error {
	if err := s.addSource(n, ctx); err != nil {
		return errors.Wrap(err, "add source")
	}

	for _, oN := range ctx.objects {
		log.Debugf("object dependency %s -> %s", n.Name, oN.Name)
		n.Dependency(oN)
	}

	for _, lN := range ctx.libraries {
		log.Debugf("library dependency %s -> %s", n.Name, lN.Name)
		n.Dependency(lN)
	}

	if ctx.cc {
		n.Resolver = s.rf.NewLinkCC()
	} else {
		n.Resolver = s.rf.NewLinkC()
	}

	return nil
}

func (s *Scanner) onObject(n *node.Node, ctx *ctx) error {
	if ctx.cc {
		n.Resolver = s.rf.NewCompileCC(ctx.includePaths)
	} else {
		n.Resolver = s.rf.NewCompileC(ctx.includePaths)
	}

	log.Debugf("adding object %s", n.Name)
	ctx.objects = append(ctx.objects, n)

	return nil
}

func (s *Scanner) onSource(n *node.Node, ctx *ctx) error {
	ext := filepath.Ext(n.Name)
	object := strings.ReplaceAll(n.Name, ext, ".o")
	objectN := node.New(object)
	if added, err := s.add(objectN, ctx); err != nil {
		return errors.Wrap(err, "add")
	} else if !added {
		return nil
	}

	log.Debugf("source/object dependency %s -> %s", objectN.Name, n.Name)
	objectN.Dependency(n)

	if err := s.addHeaders(n, ctx); err != nil {
		return errors.Wrap(err, "add headers")
	}

	return nil
}

func (s *Scanner) onHeader(n *node.Node, ctx *ctx) error {
	if err := s.addHeaders(n, ctx); err != nil {
		return errors.Wrap(err, "add headers")
	}

	if err := s.addSource(n, ctx); err != nil {
		return errors.Wrap(err, "add source")
	}

	return nil
}

func (s *Scanner) addSource(n *node.Node, ctx *ctx) error {
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
		ctx.cc = true
	}

	if sourceN == nil {
		return errors.New("unknown source for node " + n.Name)
	}

	if added, err := s.add(sourceN, ctx); err != nil {
		return errors.Wrap(err, "add "+sourceN.Name)
	} else if !added {
		return nil
	}

	return nil
}

func (s *Scanner) addHeaders(n *node.Node, ctx *ctx) error {
	includes, err := s.i.Includes(n.Name)
	if err != nil {
		return errors.Wrap(err, "includes")
	}

	for _, include := range includes {
		var headerN *node.Node
		rootRelInclude := filepath.Join(s.root, include)
		dirRelInclude := filepath.Join(filepath.Dir(n.Name), include)
		if exists, err := afero.Exists(s.fs, rootRelInclude); err != nil {
			return errors.Wrap(err, "exists")
		} else if exists {
			headerN = node.New(rootRelInclude)
		} else if exists, err = afero.Exists(s.fs, dirRelInclude); err != nil {
			return errors.Wrap(err, "exists")
		} else if exists {
			headerN = node.New(dirRelInclude)
		}

		if headerN == nil {
			return errors.New("unknown header for include " + include)
		}

		if _, err := s.add(headerN, ctx); err != nil {
			return errors.Wrap(err, "add "+headerN.Name)
		}

		log.Debugf("include dependency %s -> %s", n.Name, headerN.Name)
		n.Dependency(headerN)
	}

	return nil
}
