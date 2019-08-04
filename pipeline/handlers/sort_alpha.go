package handlers

import (
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/pipeline"
)

type sortAlpha struct {
}

// NewSortAlpha returns a pipeline.Handler that sorts a node.Node list and its
// dependencies in alphanumeric order.
func NewSortAlpha() pipeline.Handler {
	return &sortAlpha{}
}

func (sa *sortAlpha) Handle(ctx *pipeline.Ctx) {
	node.SortAlpha(ctx.Nodes)
}

func (sa *sortAlpha) Name() string { return "alpha sort" }
