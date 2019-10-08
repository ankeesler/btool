package cc

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ankeesler/btool/app/collector"
	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
)

// Exe is a collector.Consumer that adds executable node.Node's.
type Exe struct {
	rf collector.ResolverFactory
}

// NewExe creates a new Exe.
func NewExe(rf collector.ResolverFactory) *Exe {
	return &Exe{
		rf: rf,
	}
}

// Consume will listen to executable node.Node's (i.e., node.Node's with no file
// extension) and create a link target for that node.Node. It walks the source
// tree backwards to find the dependent objects and libraries.
func (e *Exe) Consume(s collector.Store, n *node.Node) error {
	if filepath.Ext(n.Name) != "" {
		return nil
	}

	c := s.Get(n.Name + ".c")
	cc := s.Get(n.Name + ".cc")
	if c != nil && cc != nil {
		return fmt.Errorf("ambiguous exe %s: %s or %s?", n, c, cc)
	} else if c == nil && cc == nil {
		return fmt.Errorf("no known source for exe %s", n)
	}

	var err error
	objs := make([]*node.Node, 0)
	if c != nil {
		log.Debugf("exe %s has source %s", n, c)
		err = collectObjs(s, c, ".c", &objs)
	} else { // cc
		log.Debugf("exe %s has source %s", n, cc)
		err = collectObjs(s, cc, ".cc", &objs)
	}

	if err != nil {
		return errors.Wrap(err, "collect objs")
	}

	libraries := make([]*node.Node, 0)
	for _, obj := range objs {
		n.Dependency(obj)

		if err := collectLibraries(s, obj, &libraries); err != nil {
			return errors.Wrap(err, "collect libraries")
		}
	}
	for _, library := range libraries {
		n.Dependency(library)
	}

	var linkFlags []string
	linkFlags, err = CollectLinkFlags(n)
	if err != nil {
		return errors.Wrap(err, "collect link flags")
	}

	var r node.Resolver
	if c != nil {
		r = e.rf.NewLinkC(linkFlags)
	} else {
		r = e.rf.NewLinkCC(linkFlags)
	}
	n.Resolver = r

	s.Set(n)

	return nil
}

// The provided n should be a c/cc file!
func collectObjs(
	s collector.Store,
	n *node.Node,
	ext string,
	objs *[]*node.Node,
) error {
	obj := s.Get(strings.ReplaceAll(n.Name, ext, ".o"))
	if obj == nil {
		return fmt.Errorf("no know object for %s", n)
	}
	if contains(*objs, obj) {
		return nil
	}
	*objs = append(*objs, obj)
	log.Debugf("remembering object %s", obj)

	for _, d := range n.Dependencies {
		if !strings.HasSuffix(d.Name, ".h") {
			continue
		}

		src := s.Get(strings.ReplaceAll(d.Name, ".h", ext))
		if src == nil {
			log.Debugf("skipping header-only dependency %s", d.Name)
			continue
		}

		if err := collectObjs(s, src, ext, objs); err != nil {
			return err
		}
	}

	return nil
}

func collectLibraries(
	s collector.Store,
	n *node.Node,
	libraries *[]*node.Node,
) error {
	if err := node.Visit(n, func(vn *node.Node) error {
		var labels Labels
		if err := collector.FromLabels(vn, &labels); err != nil {
			return errors.Wrap(err, "from labels")
		}

		for _, library := range labels.Libraries {
			libraryN := s.Get(library)
			if libraryN == nil {
				return fmt.Errorf("unknown node for library %s", library)
			}

			*libraries = append(*libraries, libraryN)
		}

		return nil
	}); err != nil {
		return errors.Wrap(err, "visit")
	}

	return nil
}

// CollectLinkFlag will walk a node.Node graph and return all of the linker flags
// paths encoutered as a part of a node.Node's Labels along the way.
func CollectLinkFlags(n *node.Node) ([]string, error) {
	return collectLabels(n, func(l *Labels) []string { return l.LinkFlags })
}

func contains(nodes []*node.Node, n *node.Node) bool {
	for _, nn := range nodes {
		if nn == n {
			return true
		}
	}
	return false
}
