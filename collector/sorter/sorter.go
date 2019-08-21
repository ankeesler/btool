// Package sorter provides a stable way of sorting a node.Node graph.
package sorter

import (
	"sort"

	"github.com/ankeesler/btool/node"
)

// Sorter is a type that can sort a node.Node graph in a stable way.
type Sorter struct {
}

// New returns a new Sorter.
func New() *Sorter {
	return &Sorter{}
}

func (s *Sorter) Sort(n *node.Node) error {
	return node.Visit(n, s.visit)
}

func (s *Sorter) visit(n *node.Node) error {
	sort.Slice(n.Dependencies, func(i, j int) bool {
		return n.Dependencies[i].Name < n.Dependencies[j].Name
	})
	return nil
}
