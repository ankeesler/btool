package handlers

import (
	"fmt"
	"io"

	"github.com/ankeesler/btool/node/pipeline"
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

func (p *print) Handle(ctx pipeline.Ctx) error {
	for _, n := range ctx.All() {
		fmt.Fprintf(p.writer, "%s\n", n.Name)
		for _, d := range n.Dependencies {
			fmt.Fprintf(p.writer, "> %s\n", d.Name)
		}
	}

	return nil
}

func (p *print) String() string { return "print" }
