package graph_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/scanner/graph"
	"github.com/sirupsen/logrus"
)

func TestSort(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

	a := &graph.Node{Name: "a"}
	b := &graph.Node{Name: "b"}
	c := &graph.Node{Name: "c"}
	d := &graph.Node{Name: "d"}
	e := &graph.Node{Name: "e"}
	f := &graph.Node{Name: "f"}

	data := []struct {
		name   string
		graph  func() *graph.Graph
		err    error
		sorted []string
	}{
		{
			name: "basic",
			graph: func() *graph.Graph {
				// a -> b -> c -> d
				//      |         ^
				//      -> e -> f /
				return graph.New().Add(a, b).Add(b, c).Add(c, d).Add(b, e).Add(e, f).Add(f, d).Add(d, nil)
			},
			err:    nil,
			sorted: []string{"d", "cf", "fc", "e", "b", "a"},
		},
		{
			name: "delayed cycle",
			graph: func() *graph.Graph {
				// a -> b -> c <-> d
				//      |         ^
				//      -> e -> f /
				return graph.New().Add(a, b).Add(b, c).Add(c, d).Add(d, c).Add(b, e).Add(e, f).Add(f, d)
			},
			err:    errors.New("cycle detected"),
			sorted: []string{},
		},
		{
			name: "immediate cycle",
			graph: func() *graph.Graph {
				// a -> b -> c
				// ^         v
				//  \ < - < /
				return graph.New().Add(a, b).Add(b, c).Add(c, a)
			},
			err:    errors.New("cycle detected"),
			sorted: []string{},
		},
		{
			name: "only dependency node (b)",
			graph: func() *graph.Graph {
				// a -> b
				return graph.New().Add(a, b)
			},
			err:    nil,
			sorted: []string{"b", "a"},
		},
	}

	for _, datum := range data {
		acNodes, err := datum.graph().Sort()
		if datum.err == nil && err != nil {
			t.Errorf("%s: expected no error, got %v", datum.name, err)
		} else if datum.err != nil {
			if err == nil {
				t.Errorf("%s: expected error, got no error", datum.name)
			} else {
				continue
			}
		}

		if ex, ac := len(datum.sorted), len(acNodes); ex != ac {
			t.Errorf("%s: expected %d nodes, got %d", datum.name, ex, ac)
		}

		for i := range acNodes {
			if ex, ac := datum.sorted[i], acNodes[i].Name; !strings.Contains(ex, ac) {
				t.Errorf("%s: node %d does not match (%s not in %s)", datum.name, i, ac, ex)
			}
		}
	}
}

func TestAdd(t *testing.T) {
	a0 := &graph.Node{Name: "a"}
	a1 := &graph.Node{Name: "a"}
	b := &graph.Node{Name: "b"}
	nodes, err := graph.New().Add(a0, nil).Add(b, nil).Add(a1, nil).Sort()
	if err != nil {
		t.Fatal(err)
	}

	l := len(nodes)
	if l != 2 {
		t.Errorf("expected length 2, got %d", l)
	}
}
