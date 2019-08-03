// Package objecter provides a node.Handler that adds objects to the nodes.
package objecter

import (
	"strings"

	"github.com/ankeesler/btool/node"
)

type Objecter struct{}

func New() *Objecter {
	return &Objecter{}
}
func (o *Objecter) Handle(cfg *node.Config, nodes []*node.Node) ([]*node.Node, error) {
	for _, n := range nodes {
		var suffix string
		if strings.HasSuffix(n.Name, ".c") {
			suffix = ".c"
		} else if strings.HasSuffix(n.Name, ".cc") {
			suffix = ".cc"
		}

		if suffix != "" {
			d := node.New(strings.ReplaceAll(n.Name, suffix, ".o"))
			d.Dependencies = append(d.Dependencies, n)
			nodes = append(nodes, d)
		}

	}
	return nodes, nil
}
