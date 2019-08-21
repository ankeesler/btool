package sorter_test

import (
	"testing"

	"github.com/ankeesler/btool/collector/sorter"
	"github.com/ankeesler/btool/node"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSorterSort(t *testing.T) {
	s := sorter.New()

	nodeA := node.New("a")
	nodeB := node.New("b")
	nodeC := node.New("c")

	node0 := node.New("0")
	node1 := node.New("1")
	node2 := node.New("2")
	node3 := node.New("3")

	nodeA.Dependency(node1, node0, nodeB)
	nodeB.Dependency(node3, node2)
	nodeC.Dependency(nodeB, nodeA)

	require.Nil(t, s.Sort(nodeA))

	assert.Equal(t, []*node.Node{node0, node1, nodeB}, nodeA.Dependencies)
	assert.Equal(t, []*node.Node{node2, node3}, nodeB.Dependencies)
	assert.Equal(t, []*node.Node{nodeB, nodeA}, nodeC.Dependencies)
}
