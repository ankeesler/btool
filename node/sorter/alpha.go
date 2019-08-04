package sorter

import (
	"sort"

	"github.com/ankeesler/btool/node"
)

// Alpha is a node.Handler that sorts node.Node's in alphanumberic order.
type Alpha struct {
}

func NewAlpha() *Alpha {
	return &Alpha{}
}

func (a *Alpha) Handle(cfg *node.Config, nodes []*node.Node) ([]*node.Node, error) {
	sorht(nodes)
	for _, n := range nodes {
		sorht(n.Dependencies)
	}
	return nodes, nil
}

func sorht(nodes []*node.Node) {
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Name < nodes[j].Name
	})
}
