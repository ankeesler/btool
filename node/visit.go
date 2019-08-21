package node

import (
	"github.com/ankeesler/btool/log"
	"github.com/pkg/errors"
)

// VisitFunc is a type of function that is called as a part of a depth-first
// traversal of a Node graph (see Visit()).
type VisitFunc func(*Node) error

// Visit performs a depth-first traversal on the provided Node graph.
// It stops and returns an error if it runs into one.
func Visit(n *Node, visitFunc VisitFunc) error {
	return visit(n, visitFunc, make(map[*Node]bool))
}

func visit(
	n *Node,
	visitFunc VisitFunc,
	visited map[*Node]bool,
) error {
	if visited[n] {
		return nil
	}

	log.Debugf("visiting %s", n.Name)

	for _, dN := range n.Dependencies {
		if err := visit(dN, visitFunc, visited); err != nil {
			return errors.Wrap(err, "visit")
		}
	}

	if err := visitFunc(n); err != nil {
		return errors.Wrap(err, "visit func")
	}

	visited[n] = true

	return nil
}
