// Package objector provide node.Handler's that create build targets in
// a node list.
package objecter

import (
	"fmt"
	"strings"

	"github.com/ankeesler/btool/node"
)

// Objecter is a node.Handler that creates build targets from a provided target.
type Objecter struct {
}

func New() *Objecter {
	return &Objecter{}
}

func (o *Objecter) Handle(cfg *node.Config, nodes []*node.Node) ([]*node.Node, error) {
	if !strings.HasSuffix(cfg.Target, ".o") {
		return nodes, nil
	}

	var d *node.Node
	c := node.Find(strings.ReplaceAll(cfg.Target, ".o", ".c"), nodes)
	cc := node.Find(strings.ReplaceAll(cfg.Target, ".o", ".cc"), nodes)
	if c != nil && cc != nil {
		return nil, fmt.Errorf("ambiguous object %s (%s or %s)", cfg.Target, c.Name, cc.Name)
	} else if c != nil {
		d = c
	} else if cc != nil {
		d = cc
	}

	if d == nil {
		return nil, fmt.Errorf("unknown source for object %s", cfg.Target)
	}

	nodes = append(nodes, node.New(cfg.Target).Dependency(d))

	return nodes, nil
}