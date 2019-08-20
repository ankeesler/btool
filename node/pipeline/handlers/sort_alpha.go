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

func (sa *sortAlpha) Handle(ctx pipeline.Ctx) error {
	// TODO: this won't work!
	node.SortAlpha(ctx.All())
	return nil
}

func (sa *sortAlpha) String() string { return "alpha sort" }
