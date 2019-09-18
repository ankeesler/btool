// Package testutil provides test utilities for the collector framework.
package testutil

import (
	"github.com/ankeesler/btool/app/collector/collectorfakes"
	"github.com/ankeesler/btool/node"
)

func FakeStore(initNodes ...*node.Node) *collectorfakes.FakeStore {
	nodes := make(map[string]*node.Node)
	for _, n := range initNodes {
		nodes[n.Name] = n
	}

	s := &collectorfakes.FakeStore{}
	s.GetStub = func(name string) *node.Node {
		return nodes[name]
		return nil
	}
	s.ForEachStub = func(f func(*node.Node)) {
		for _, n := range nodes {
			f(n)
		}
	}
	s.SetStub = func(n *node.Node) {
		nodes[n.Name] = n
	}

	return s
}
