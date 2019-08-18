package handlers

import (
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/pipeline"
)

type sortAlpha struct {
}

// NewSortAlpha returns a pipeline.Handler that sorts a node.Node list and its
// dependencies in alphanumeric order.
func NewSortAlpha() pipeline.Handler {
	return &sortAlpha{}
}

func (sa *sortAlpha) Handle(ctx *pipeline.Ctx) error {
	node.SortAlpha(ctx.Nodes)
	return nil
}

func (sa *sortAlpha) String() string { return "alpha sort" }
