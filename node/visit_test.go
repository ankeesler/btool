package node_test

import (
	"testing"

	"github.com/ankeesler/btool/node"
	"github.com/stretchr/testify/require"
)

func TestVisit(t *testing.T) {
	// a -> b, c
	// b -> c
	// c
	nodeC := node.New("c")
	nodeB := node.New("b").Dependency(nodeC)
	nodeA := node.New("a").Dependency(nodeB, nodeC)

	data := []struct {
		name      string
		n         *node.Node
		exVisited []*node.Node
	}{
		{
			name:      "All",
			n:         nodeA,
			exVisited: []*node.Node{nodeC, nodeB, nodeA},
		},
	}
	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			acVisited := make([]*node.Node, 0)
			visitFunc := func(n *node.Node) error {
				acVisited = append(acVisited, n)
				return nil
			}

			require.Nil(t, node.Visit(datum.n, visitFunc))
			require.Equal(t, datum.exVisited, acVisited)
		})
	}
}
