package handlers

import (
	"fmt"
	"io"

	"github.com/ankeesler/btool/pipeline"
)

type print struct {
	writer io.Writer
}

// NewPrint returns a pipeline.Handler that prints out the current Ctx to
// the provided io.Writer.
func NewPrint(writer io.Writer) pipeline.Handler {
	return &print{
		writer: writer,
	}
}

func (p *print) Handle(ctx *pipeline.Ctx) {
	fmt.Fprintln(p.writer, "*** Nodes ***")
	for _, n := range ctx.Nodes {
		fmt.Fprintf(p.writer, "%s\n", n.Name)
		for _, d := range n.Dependencies {
			fmt.Fprintf(p.writer, "> %s\n", d.Name)
		}
	}

	fmt.Fprintln(p.writer, "*** Err ***")
	fmt.Fprintln(p.writer, ctx.Err)

	fmt.Fprintln(p.writer, "*** KV ***")
	fmt.Fprintln(p.writer, ctx.KV)
}

func (p *print) Name() string { return "print" }
