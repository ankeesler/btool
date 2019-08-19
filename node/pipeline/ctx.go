package pipeline

import "github.com/ankeesler/btool/node"

type ctx struct {
	cb    Callback
	nodes []*node.Node
}

func newCtx(cb Callback) *ctx {
	return &ctx{
		cb:    cb,
		nodes: make([]*node.Node, 0),
	}
}

func (c *ctx) Add(n *node.Node) {
	c.cb.OnAdd(n)
	c.nodes = append(c.nodes, n)
}

func (c *ctx) Find(name string) *node.Node {
	for _, n := range c.nodes {
		if n.Name == name {
			return n
		}
	}
	return nil
}

func (c *ctx) All() []*node.Node {
	return c.nodes
}
