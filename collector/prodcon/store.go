package prodcon

import "github.com/ankeesler/btool/node"

// Store is a place where node.Node's are kept.
type Store struct {
	nodes map[string]*node.Node

	addCallback func(*node.Node)
}

func NewStore() *Store {
	return &Store{
		nodes: make(map[string]*node.Node),
	}
}

func (s *Store) Add(nodes ...*node.Node) {
	for _, n := range nodes {
		s.nodes[n.Name] = n
		if s.addCallback != nil {
			s.addCallback(n)
		}
		//if s.w != nil {
		//	s.w.OnAdd(n)
		//}
	}
}

func (s *Store) Find(name string) *node.Node {
	return s.nodes[name]
}
