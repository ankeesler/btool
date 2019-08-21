package builder_test

import (
	"testing"

	"github.com/ankeesler/btool/builder"
	"github.com/ankeesler/btool/builder/builderfakes"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/nodefakes"
	"github.com/stretchr/testify/assert"
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
		current    []string
		exResolved []string
		exCallback []string
	}{
		{
			name:       "All",
			nodes:      nodes,
			target:     "a",
			current:    []string{},
			exResolved: []string{"c", "b", "a"},
			exCallback: []string{"c", "b", "a"},
		},
		{
			name:       "UpToDate",
			nodes:      nodes,
			target:     "a",
			current:    []string{"a", "b", "c"},
			exResolved: []string{},
			exCallback: []string{"c", "b", "a"},
		},
		{
			name:       "Some",
			nodes:      nodes,
			target:     "a",
			current:    []string{"c", "b"},
			exResolved: []string{"a"},
			exCallback: []string{"c", "b", "a"},
		},
	}
	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			c := &builderfakes.FakeCurrenter{}
			c.CurrentStub = func(n *node.Node) (bool, error) {
				return contains(n.Name, datum.current), nil
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

			n := find(datum.target, datum.nodes)
			require.NotNil(t, n)

			callback := &builderfakes.FakeCallback{}

			b := builder.New(false, c, callback)
			require.Nil(t, b.Build(n))
			require.Equal(t, datum.exResolved, acResolved)

			assert.Equal(t, len(datum.exCallback), callback.OnResolveCallCount())
			for i, exName := range datum.exCallback {
				exCurrent := contains(exName, datum.current)
				exN := find(exName, datum.nodes)
				require.NotNil(t, exN)

				acN, acCurrent := callback.OnResolveArgsForCall(i)

				assert.Equal(t, exN, acN)
				assert.Equal(t, exCurrent, acCurrent)
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
