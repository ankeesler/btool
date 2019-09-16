package cc

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ankeesler/btool/collector"
	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Includeser

// Includeser is a type that can return a list of #include's from a given file.
type Includeser interface {
	Includes(path string) ([]string, error)
}

type Includes struct {
	i Includeser
}

func NewIncludes(i Includeser) *Includes {
	return &Includes{
		i: i,
	}
}

func (i *Includes) Consume(s collector.Store, n *node.Node) error {
	ext := filepath.Ext(n.Name)
	if ext != ".c" && ext != ".cc" && ext != ".h" {
		return nil
	}

	// TODO: another string conversation mechanism needed here.
	if l, ok := n.Labels[collector.LabelLocal]; !ok || l != "true" {
		return nil
	}

	includes, err := i.i.Includes(n.Name)
	if err != nil {
		return errors.Wrap(err, "includes")
	}

	for _, include := range includes {
		d, err := i.resolveInclude(s, n, include)
		if err != nil {
			return errors.Wrap(err, "resolve include")
		}
		n.Dependency(d)
		log.Debugf("include dependency %s -> %s", n, d)
	}

	s.Set(n)

	return nil
}

func (i *Includes) resolveInclude(
	s collector.Store,
	n *node.Node,
	include string,
) (*node.Node, error) {
	dirLocal := filepath.Join(filepath.Dir(n.Name), include)
	d := s.Get(dirLocal)
	log.Debugf(
		"checking dir local %s for include %s for node %s",
		dirLocal,
		include,
		n,
	)
	if d != nil {
		return d, nil
	}

	// TODO: this is ugly! We could pick up the wrong include based on suffix!
	var includePath string
	s.ForEach(func(sn *node.Node) {
		if d != nil {
			return
		}

		log.Debugf("does node %s have suffix %s", sn.Name, include)
		if !strings.HasSuffix(sn.Name, include) {
			return
		}

		includePath = strings.TrimSuffix(sn.Name, include)
		if includePath == "" {
			includePath = "."
		}
		log.Debugf("yes, and include path is %s", includePath)

		d = sn
		n.Labels[LabelIncludePaths] += includePath + ","
	})
	if d != nil {
		return d, nil
	}

	return nil, fmt.Errorf("cannot resolve include %s for node %s", include, n)
}
