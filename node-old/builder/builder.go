package builder

import "github.com/ankeesler/btool/node"

type Builder struct {
}

func New() *Builder {
	return &Builder{}
}

func (b *Builder) Handle(c *node.Config, nodes []*node.Node) ([]*node.Node, error) {
	return nodes, nil
}
