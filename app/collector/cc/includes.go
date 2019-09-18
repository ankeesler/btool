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

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Includeser

// Includeser is a type that can return a list of #include's from a given file.
type Includeser interface {
	Includes(path string) ([]string, error)
}

// Includes is a collector.Consumer that adds dependencies to local C/C++
// node.Node's.
type Includes struct {
	i Includeser
}

// NewIncludes creates a new Includes.
func NewIncludes(i Includeser) *Includes {
	return &Includes{
		i: i,
	}
}

// Consume will react to local .c/.cc/.h files and add Dependencies and
// IncludePaths to a node.Node depending on what #include's the file contains.
func (i *Includes) Consume(s collector.Store, n *node.Node) error {
	ext := filepath.Ext(n.Name)
	if ext != ".c" && ext != ".cc" && ext != ".h" {
		return nil
	}

	if local, err := isLocal(n); err != nil {
		return errors.Wrap(err, "is local")
	} else if !local {
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
	var err error
	s.ForEach(func(sn *node.Node) {
		if d != nil {
			return
		}

		log.Debugf("does node %s have suffix %s", sn.Name, include)
		if !strings.HasSuffix(sn.Name, include) {
			return
		}

		includePath := strings.TrimSuffix(sn.Name, include)
		if includePath == "" {
			includePath = "."
		}
		log.Debugf("yes, and include path is %s", includePath)

		d = sn
		err = AppendIncludePaths(n, includePath)
	})
	if err != nil {
		return nil, errors.Wrap(err, "for each")
	} else if d != nil {
		return d, nil
	}

	return nil, fmt.Errorf("cannot resolve include %s for node %s", include, n)
}

func isLocal(n *node.Node) (bool, error) {
	var labels collector.Labels
	if err := collector.FromLabels(n, &labels); err != nil {
		return false, errors.Wrap(err, "from labels")
	}

	return labels.Local, nil
}

// AppendIncludePaths will append the provided include paths to the node.Node's
// Labels.
func AppendIncludePaths(n *node.Node, includePaths ...string) error {
	var labels Labels
	if err := collector.FromLabels(n, &labels); err != nil {
		return errors.Wrap(err, "from labels")
	}

	labels.IncludePaths = append(labels.IncludePaths, includePaths...)

	if err := collector.ToLabels(n, &labels); err != nil {
		return errors.Wrap(err, "to labels")
	}

	return nil
}
