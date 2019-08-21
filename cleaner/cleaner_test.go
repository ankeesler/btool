package cleaner_test

import (
	"testing"

	"github.com/ankeesler/btool/cleaner"
	"github.com/ankeesler/btool/cleaner/cleanerfakes"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/nodefakes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCleanerClean(t *testing.T) {
	// a -> b, c
	// b -> c
	// c
	// d -> e
	// e -> f
	// f
	nodeF := node.New("f")
	nodeE := node.New("e").Dependency(nodeF)
	nodeD := node.New("d").Dependency(nodeE)
	nodeC := node.New("c")
	nodeB := node.New("b").Dependency(nodeC)
	nodeA := node.New("a").Dependency(nodeB, nodeC)

	nodeB.Resolver = &nodefakes.FakeResolver{}
	nodeA.Resolver = &nodefakes.FakeResolver{}

	nodes := []*node.Node{
		nodeA,
		nodeB,
		nodeC,
		nodeD,
		nodeE,
		nodeF,
	}

	data := []struct {
		name       string
		nodes      []*node.Node
		target     string
		exCleaned  []string
		exCallback []string
	}{
		{
			name:      "All",
			nodes:     nodes,
			target:    "a",
			exCleaned: []string{"b", "a"},
		},
	}
	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			ra := &cleanerfakes.FakeRemoveAller{}
			cb := &cleanerfakes.FakeCallback{}
			c := cleaner.New(ra, cb)

			n := find(datum.target, datum.nodes)
			require.NotNil(t, n)
			require.Nil(t, c.Clean(n))

			assert.Equal(t, len(datum.exCleaned), ra.RemoveAllCallCount())
			assert.Equal(t, len(datum.exCleaned), cb.OnCleanCallCount())
			for i, exName := range datum.exCleaned {
				acName := ra.RemoveAllArgsForCall(i)
				assert.Equal(t, exName, acName)

				exN := find(exName, datum.nodes)
				require.NotNil(t, exN)

				acN := cb.OnCleanArgsForCall(i)
				assert.Equal(t, exN, acN)
			}
		})
	}
}

func contains(s string, ss []string) bool {
	for _, tuna := range ss {
		if tuna == s {
			return true
		}
	}
	return false
}

func find(target string, nodes []*node.Node) *node.Node {
	for _, n := range nodes {
		if n.Name == target {
			return n
		}
	}
	return nil
}
