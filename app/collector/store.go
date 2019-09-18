package collector

import "github.com/ankeesler/btool/node"

type store struct {
	nodes map[string]*node.Node

	setCallback func(*node.Node)
}

func newStore() *store {
	return &store{
		nodes: make(map[string]*node.Node),
	}
}

func (s *store) Set(n *node.Node) {
	s.nodes[n.Name] = n
	if s.setCallback != nil {
		s.setCallback(n)
	}
}

func (s *store) Get(name string) *node.Node {
	return s.nodes[name]
}

func (s *store) ForEach(f func(*node.Node)) {
	for _, n := range s.nodes {
		f(n)
	}
}
