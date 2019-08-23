// Package nodestore provides a type that can create and find node.Node's.
package nodestore

import "github.com/ankeesler/btool/node"

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Watcher

// A Watcher can be notified of when node.Node's are added to the NodeStore.
type Watcher interface {
	OnAdd(*node.Node)
}

// NodeStore is a type that can create a find node.Node's.
type NodeStore struct {
	nodes map[string]*node.Node
	w     Watcher
}

// New creates a new NodeStore with a Watcher.
//
// The provided Watcher can be nil.
func New(w Watcher) *NodeStore {
	return &NodeStore{
		nodes: make(map[string]*node.Node),
		w:     w,
	}
}

func (ns *NodeStore) Add(n *node.Node) {
	ns.nodes[n.Name] = n
	if ns.w != nil {
		ns.w.OnAdd(n)
	}
}

func (ns *NodeStore) Find(name string) *node.Node {
	return ns.nodes[name]
}
