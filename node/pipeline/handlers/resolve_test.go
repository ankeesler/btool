package handlers_test

import (
	"testing"
	"time"

	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/nodefakes"
	"github.com/ankeesler/btool/node/pipeline"
	"github.com/ankeesler/btool/node/pipeline/handlers"
	"github.com/go-test/deep"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func TestResolve(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

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
		written    []string
		exResolved []string
	}{
		{
			name:       "All",
			nodes:      nodes,
			target:     "a",
			written:    []string{},
			exResolved: []string{"c", "b", "a"},
		},
		{
			name:       "UpToDate",
			nodes:      nodes,
			target:     "a",
			written:    []string{"c", "b", "a"},
			exResolved: []string{},
		},
		{
			name:       "Some",
			nodes:      nodes,
			target:     "a",
			written:    []string{"c", "b"},
			exResolved: []string{"a"},
		},
		{
			name:       "Newer",
			nodes:      nodes,
			target:     "a",
			written:    []string{"a", "c"},
			exResolved: []string{"b", "a"},
		},
		{
			name:       "DidNotExist",
			nodes:      nodes,
			target:     "a",
			written:    []string{"b", "a"},
			exResolved: []string{"c", "b", "a"},
		},
		{
			name:       "TwoDeep",
			nodes:      nodes,
			target:     "d",
			written:    []string{"d", "e"},
			exResolved: []string{"f", "e", "d"},
		},
	}
	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()

			acResolved := make([]string, 0)
			for _, n := range nodes {
				r := &nodefakes.FakeResolver{}
				r.ResolveStub = func(n *node.Node) error {
					if err := afero.WriteFile(
						fs,
						n.Name,
						[]byte(n.Name),
						0644,
					); err != nil {
						t.Fatal(err)
					}
					acResolved = append(acResolved, n.Name)
					return nil
				}
				n.Resolver = r
			}

			for _, w := range datum.written {
				time.Sleep(time.Millisecond * 250)
				if err := afero.WriteFile(fs, w, []byte(w), 0644); err != nil {
					t.Fatal(err)
				}
			}

			h := handlers.NewResolve(fs)
			ctx := pipeline.NewCtxBuilder().Nodes(
				datum.nodes,
			).Target(
				datum.target,
			).Build()
			if err := h.Handle(ctx); err != nil {
				t.Fatal(err)
			}

			if diff := deep.Equal(datum.exResolved, acResolved); diff != nil {
				t.Fatal(diff)
			}
		})
	}
}
