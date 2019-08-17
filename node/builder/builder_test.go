package builder_test

import (
	"testing"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/builder"
	"github.com/ankeesler/btool/node/builder/builderfakes"
	"github.com/ankeesler/btool/node/nodefakes"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/require"
)

func TestBuilderBuild(t *testing.T) {
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
		current    map[string]bool
		exResolved []string
	}{
		{
			name:   "All",
			nodes:  nodes,
			target: "a",
			current: map[string]bool{
				"a": false,
				"b": false,
				"c": false,
			},
			exResolved: []string{"c", "b", "a"},
		},
		{
			name:   "UpToDate",
			nodes:  nodes,
			target: "a",
			current: map[string]bool{
				"a": true,
				"b": true,
				"c": true,
			},
			exResolved: []string{},
		},
		{
			name:   "Some",
			nodes:  nodes,
			target: "a",
			current: map[string]bool{
				"a": false,
				"b": true,
				"c": true,
			},
			exResolved: []string{"a"},
		},
	}
	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			c := &builderfakes.FakeCurrenter{}
			c.CurrentStub = func(n *node.Node) (bool, error) {
				current := datum.current[n.Name]
				return current, nil
			}

			acResolved := make([]string, 0)
			for _, n := range nodes {
				r := &nodefakes.FakeResolver{}
				r.ResolveStub = func(n *node.Node) error {
					acResolved = append(acResolved, n.Name)
					return nil
				}
				n.Resolver = r
			}

			n := node.Find(datum.target, datum.nodes)
			require.NotNil(t, n)

			b := builder.New(c)
			require.Nil(t, b.Build(n))
			require.Nil(t, deep.Equal(datum.exResolved, acResolved))
		})
	}
}
