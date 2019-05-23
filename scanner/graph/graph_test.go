package graph_test

import (
	"errors"
	"reflect"
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
			sorted: []string{"d", "c", "f", "e", "b", "a"},
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

func TestWalk(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

	a := &graph.Node{Name: "a"}
	b := &graph.Node{Name: "b"}
	c := &graph.Node{Name: "c"}
	d := &graph.Node{Name: "d"}
	e := &graph.Node{Name: "e"}

	data := []struct {
		name    string
		graph   func() *graph.Graph
		root    string
		visits  []string
		failure bool
	}{
		{
			name: "Basic1",
			graph: func() *graph.Graph {
				// a -> b -> c
				// d -> e
				return graph.New().Add(a, b).Add(b, c).Add(d, e)
			},
			root:    "a",
			visits:  []string{"a", "b", "c"},
			failure: false,
		},
		{
			name: "Basic2",
			graph: func() *graph.Graph {
				// a -> b -> c
				// d -> e
				return graph.New().Add(a, b).Add(b, c).Add(d, e)
			},
			root:    "d",
			visits:  []string{"d", "e"},
			failure: false,
		},
		{
			name: "EndOfAPath",
			graph: func() *graph.Graph {
				// a -> b -> c
				// d -> e
				return graph.New().Add(a, b).Add(b, c).Add(d, e)
			},
			root:    "e",
			visits:  []string{"e"},
			failure: false,
		},
		{
			name: "Cycle",
			graph: func() *graph.Graph {
				// a -> b -> c
				// ^         v
				//  \ < - < /
				return graph.New().Add(a, b).Add(b, c).Add(c, a)
			},
			root:    "a",
			visits:  []string{"a", "b", "c"},
			failure: false,
		},
		{
			name: "UnknownNode",
			graph: func() *graph.Graph {
				// a -> b -> c
				return graph.New().Add(a, b).Add(b, c)
			},
			root:    "z",
			visits:  []string{},
			failure: true,
		},
	}

	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			g := datum.graph()

			acVisits := make([]string, 0, 2)
			if err := g.Walk(datum.root, func(node *graph.Node) error {
				acVisits = append(acVisits, node.Name)
				return nil
			}); err != nil && !datum.failure {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(datum.visits, acVisits) {
				t.Fatalf("expected visits %s != actual visits %s", datum.visits, acVisits)
			}
		})
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

func TestString(t *testing.T) {
	a := &graph.Node{Name: "a"}
	b := &graph.Node{Name: "b"}
	c := &graph.Node{Name: "c"}
	makeGraph := func() *graph.Graph {
		// a -> b -> c
		return graph.New().Add(a, b).Add(b, c).Add(c, a)
	}

	g0, g1 := makeGraph(), makeGraph()
	g0S, g1S := g0.String(), g1.String()
	if g0S != g1S {
		t.Fatalf("'%s' != '%s'", g0S, g1S)
	}
}

func TestEqual(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

	a := &graph.Node{Name: "a"}
	b := &graph.Node{Name: "b"}
	c := &graph.Node{Name: "c"}
	d := &graph.Node{Name: "d"}
	e := &graph.Node{Name: "e"}

	data := []struct {
		name    string
		g0, g1  *graph.Graph
		failure bool
	}{
		{
			// a -> b -> c
			// vs
			// a -> b -> c
			name:    "Correct",
			g0:      graph.New().Add(a, b).Add(b, c),
			g1:      graph.New().Add(a, b).Add(b, c),
			failure: false,
		},
		{
			// a -> b -> c
			// vs
			// d -> e
			name:    "Wrong",
			g0:      graph.New().Add(a, b).Add(b, c),
			g1:      graph.New().Add(d, e),
			failure: true,
		},
		{
			// a -> b -> c
			// ^         v
			//  \ < - < /
			// vs
			// a -> b -> c
			// ^         v
			//  \ < - < /
			name:    "Cycle",
			g0:      graph.New().Add(a, b).Add(b, c).Add(c, a),
			g1:      graph.New().Add(a, b).Add(b, c).Add(c, a),
			failure: false,
		},
	}

	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			err := graph.Equal(datum.g0, datum.g1)
			if err != nil {
				if !datum.failure {
					t.Fatalf("expected success, got %s", err.Error())
				}
			} else {
				if datum.failure {
					t.Fatal("expected failure")
				}
			}
		})
	}
}
