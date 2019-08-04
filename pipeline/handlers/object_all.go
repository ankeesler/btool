package handlers

import (
	"strings"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/pipeline"
)

type objectAll struct{}

// NewObjectAll returns a pipeline.Handler that creates objects for each .c/.cc
// node.Node in the node.Node list.
func NewObjectAll() pipeline.Handler {
	return &objectAll{}
}

func (oa *objectAll) Handle(ctx *pipeline.Ctx) {
	for _, n := range ctx.Nodes {
		var suffix string
		if strings.HasSuffix(n.Name, ".c") {
			suffix = ".c"
		} else if strings.HasSuffix(n.Name, ".cc") {
			suffix = ".cc"
		}

		if suffix != "" {
			d := node.New(strings.ReplaceAll(n.Name, suffix, ".o"))
			d.Dependencies = append(d.Dependencies, n)
			ctx.Nodes = append(ctx.Nodes, d)
		}

	}
}

func (oa *objectAll) Name() string { return "all objects" }
