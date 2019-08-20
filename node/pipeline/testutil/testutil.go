package testutil

import "github.com/ankeesler/btool/node"

type Ctx struct {
	//cb    Callback
	Nodes []*node.Node
}

func NewCtx() *Ctx {
	return &Ctx{
		//cb:    cb,
		Nodes: make([]*node.Node, 0),
	}
}

func (c *Ctx) Add(n *node.Node) {
	//c.cb.OnAdd(n)
	c.Nodes = append(c.Nodes, n)
}

func (c *Ctx) Find(name string) *node.Node {
	for _, n := range c.Nodes {
		if n.Name == name {
			return n
		}
	}
	return nil
}

func (c *Ctx) All() []*node.Node {
	return c.Nodes
}
