package cc

import (
	"fmt"
	"path/filepath"
	"strings"

	collector "github.com/ankeesler/btool/collector0"
	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
)

type Exe struct {
	rf ResolverFactory
}

func NewExe(rf ResolverFactory) *Exe {
	return &Exe{
		rf: rf,
	}
}

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

	var r node.Resolver
	var err error
	objs := make(map[string]*node.Node)
	if c != nil {
		log.Debugf("exe %s has source %s", n, c)

		err = collectObjs(s, c, ".c", objs)
		r = e.rf.NewLinkC()
	} else { // cc
		log.Debugf("exe %s has source %s", n, cc)

		err = collectObjs(s, cc, ".cc", objs)
		r = e.rf.NewLinkCC()
	}

	if err != nil {
		return errors.Wrap(err, "collect objs")
	}

	for _, obj := range objs {
		n.Dependency(obj)
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
	objs map[string]*node.Node,
) error {
	obj := s.Get(strings.ReplaceAll(n.Name, ext, ".o"))
	if obj == nil {
		return fmt.Errorf("no know object for %s", n)
	}
	if _, ok := objs[obj.Name]; ok {
		return nil
	}
	objs[obj.Name] = obj
	log.Debugf("remembering object %s", obj)

	for _, d := range n.Dependencies {
		if !strings.HasSuffix(d.Name, ".h") {
			continue
		}

		src := s.Get(strings.ReplaceAll(d.Name, ".h", ext))
		if err := collectObjs(s, src, ext, objs); err != nil {
			return err
		}
	}

	return nil
}
