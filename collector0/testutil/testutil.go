// Package testutil provides test utilities for the collector framework.
package testutil

import (
	"github.com/ankeesler/btool/collector0/collector0fakes"
	"github.com/ankeesler/btool/node"
)

func FakeStore(nodes ...*node.Node) *collector0fakes.FakeStore {
	s := &collector0fakes.FakeStore{}
	s.GetStub = func(name string) *node.Node {
		for _, n := range nodes {
			if n.Name == name {
				return n
			}
		}
		return nil
	}
	s.ForEachStub = func(f func(*node.Node)) {
		for _, n := range nodes {
			f(n)
		}
	}
	return s
}
