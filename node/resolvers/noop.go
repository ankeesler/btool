package resolvers

import (
	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
)

type noop struct {
}

// NewNoop provides a node.Resolver that simply prints out that it was called.
func NewNoop() node.Resolver {
	return &noop{}
}

func (noop *noop) Resolve(n *node.Node) error {
	log.Infof("resolve %s", n)
	return nil
}
