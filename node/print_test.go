package node_test

import (
	"testing"

	"github.com/ankeesler/btool/node"
	"github.com/stretchr/testify/assert"
)

func TestPrint(t *testing.T) {
	// a -> b, c
	// b -> c
	// c -> d
	// d -> e
	// e -> f
	// f
	nodeF := node.New("f")
	nodeE := node.New("e").Dependency(nodeF)
	nodeD := node.New("d").Dependency(nodeE)
	nodeC := node.New("c").Dependency(nodeD)
	nodeB := node.New("b").Dependency(nodeC)
	nodeA := node.New("a").Dependency(nodeB, nodeC, nodeF)

	data := []struct {
		n  *node.Node
		ex string
	}{
		{
			n: nodeA,
			ex: `a
. b
. . c
. . . d
. . . . e
. . . . . f
. c
. . d
. . . e
. . . . f
. f
`,
		},
	}
	for _, datum := range data {
		ac := node.String(datum.n)
		assert.Equal(t, datum.ex, ac)
	}
}
