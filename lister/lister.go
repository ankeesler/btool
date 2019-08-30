// Package lister provides a type that can list all of the node.Node's in a
// node.Node graph.
package lister

import (
	"fmt"
	"io"

	"github.com/ankeesler/btool/node"
)

// Lister is a type that can list all of the node.Node's in a node.Node graph.
type Lister struct {
	w io.Writer
}

// New creates a new Lister.
func New(w io.Writer) *Lister {
	return &Lister{
		w: w,
	}
}

func (l *Lister) List(n *node.Node) error {
	fmt.Fprintln(l.w, node.String(n))
	return nil
}
