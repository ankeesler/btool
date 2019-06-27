// Package printer provides a node.Hander that prints nodes.
package printer

import (
	"fmt"
	"io"

	"github.com/ankeesler/btool/node"
)

type Printer struct {
	writer io.Writer
}

func New(writer io.Writer) *Printer {
	return &Printer{
		writer: writer,
	}
}

func (p *Printer) Handle(nodes []*node.Node) ([]*node.Node, error) {
	for _, node := range nodes {
		fmt.Fprintf(p.writer, "%s\n", node)
		for _, dependency := range node.Dependencies {
			fmt.Fprintf(p.writer, "> %s\n", dependency)
		}
	}
	return nodes, nil
}
