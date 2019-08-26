// Package collector provides functionality to build a node.Node graph.
package collector

import (
	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Collectini

// Collectini is an object that contributes some node.Node's to a graph.
// It should return an error if it fails to Collect() all the node.Node's that it
// cares about.
type Collectini interface {
	Collect(*Ctx, *node.Node) error
}

// Collector is a type that can build a node.Node graph. It does this via
// Collectini's that each provide a different part of the node.Node graph.
type Collector struct {
	ctx    *Ctx
	ctinis []Collectini
}

// New creates a new Collector.
func New(
	ctx *Ctx,
	ctinis ...Collectini,
) *Collector {
	return &Collector{
		ctx:    ctx,
		ctinis: ctinis,
	}
}

// Collect creates a node.Node graph.
func (c *Collector) Collect(n *node.Node) error {
	for _, ctini := range c.ctinis {
		if err := ctini.Collect(c.ctx, n); err != nil {
			return errors.Wrap(err, "collect")
		}
	}
	return nil
}
