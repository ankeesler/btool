// Package compiler provides a node.Handler that compiles each node into an object
// archive.
package compiler

import (
	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . C

// C is an actual compiler client.
type C interface {
	CompileC(output, input string, includeDirs []string) error
	CompileCC(output, input string, includeDirs []string) error
}

type Compiler struct {
	c     C
	root  string
	cache string
}

func New(c C, root, cache string) *Compiler {
	return &Compiler{
		c:     c,
		root:  root,
		cache: cache,
	}
}

func (c *Compiler) Handle(nodes []*node.Node) ([]*node.Node, error) {
	for _, node := range nodes {
		if err := c.handleNode(node); err != nil {
			return nil, errors.Wrap(err, "handle node "+node.Name)
		}
	}
	return nodes, nil
}

func (c *Compiler) handleNode(n *node.Node) error {
	if !c.needsCompile(n) {
		return nil
	}

	for _, source := range n.Sources {
		if err := c.handleSource(source, n.IncludePaths); err != nil {
			return errors.Wrap(err, "handle source "+source)
		}
	}
	return nil
}

func (c *Compiler) needsCompile(n *node.Node) bool {
	return true
}

func (c *Compiler) handleSource(source string, includePaths []string) error {
	// Get correct compiler func.
	// Create directory in cache based on source file directory.
	// Create filepath for object.
	// Invoke compiler with output, input, include directories
	return nil
}
