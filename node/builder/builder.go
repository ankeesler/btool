// Package builder provides a node.Handler that actually runs the node.Resolver's
// in the node list.
package builder

import (
	"fmt"

	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Builder struct{}

func New() *Builder {
	return &Builder{}
}

func (b *Builder) Handle(cfg *node.Config, nodes []*node.Node) ([]*node.Node, error) {
	n := node.Find(cfg.Target, nodes)
	if n == nil {
		return nil, fmt.Errorf("unknown target %s", cfg.Target)
	}

	if err := build(n); err != nil {
		return nil, errors.Wrap(err, "build")
	}

	return nodes, nil
}

func build(n *node.Node) error {
	logrus.Debugf("building %s", n.Name)

	for _, d := range n.Dependencies {
		logrus.Debugf("building dependency %s", d.Name)
		if err := build(d); err != nil {
			return errors.Wrap(err, "build "+d.Name)
		}
	}

	logrus.Debugf("resolving %s", n.Name)
	if err := n.Resolver.Resolve(n); err != nil {
		return errors.Wrap(err, "resolve "+n.Name)
	}

	return nil
}
