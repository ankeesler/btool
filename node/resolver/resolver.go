// Package resolver provides various node.Resolver's to a node.Node list.
package resolver

import "github.com/ankeesler/btool/node"

type Resolver struct {
}

func New() *Resolver {
	return &Resolver{}
}

func (r *Resolver) Handle(cfg *node.Config, nodes []*node.Node) ([]*node.Node, error) {
	return nodes, nil
}
