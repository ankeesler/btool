package handlers_test

import (
	"testing"

	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/nodefakes"
	"github.com/ankeesler/btool/pipeline"
	"github.com/ankeesler/btool/pipeline/handlers"
	"github.com/go-test/deep"
	"github.com/sirupsen/logrus"
)

func TestResolve(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

	// a -> b, c
	// b -> c
	// c
	nodeC := node.New("c")
	nodeB := node.New("b").Dependency(nodeC)
	nodeA := node.New("a").Dependency(nodeB, nodeC)

	nodes := []*node.Node{
		nodeA,
		nodeB,
		nodeC,
	}

	data := []struct {
		name       string
		nodes      []*node.Node
		target     string
		exResolved []string
	}{
		{
			name:       "All",
			nodes:      nodes,
			target:     "a",
			exResolved: []string{"c", "b", "a"},
		},
	}
	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			acResolved := make([]string, 0)
			for _, n := range nodes {
				r := &nodefakes.FakeResolver{}
				r.ResolveStub = func(n *node.Node) error {
					acResolved = append(acResolved, n.Name)
					return nil
				}
				n.Resolver = r
			}

			h := handlers.NewResolve()
			ctx := pipeline.NewCtxBuilder().Nodes(
				datum.nodes,
			).Target(
				datum.target,
			).Build()
			h.Handle(ctx)
			if ctx.Err != nil {
				t.Fatal(ctx.Err)
			}

			if diff := deep.Equal(datum.exResolved, acResolved); diff != nil {
				t.Fatal(diff)
			}
		})
	}
}
