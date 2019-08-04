package handlers

import (
	"sort"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/pipeline"
)

type sortAlpha struct {
}

// NewSortAlpha returns a pipeline.Handler that sorts a node.Node list and its
// dependencies in alphanumeric order.
func NewSortAlpha() *sortAlpha {
	return &sortAlpha{}
}

func (sa *sortAlpha) Handle(ctx *pipeline.Ctx) {
	sorht(ctx.Nodes)
	for _, n := range ctx.Nodes {
		sorht(n.Dependencies)
	}
}

func sorht(nodes []*node.Node) {
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Name < nodes[j].Name
	})
}

func (sa *sortAlpha) Name() string { return "alpha sort" }
