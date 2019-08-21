// Package collector provides functionality to build a node.Node graph.
package collector

import (
	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Scanner

// Scanner is an object that can build up a local node.Node graph.
type Scanner interface {
	// Given a starting node.Node, the Scanner should walk the node.Node graph
	// in the local FS.
	Scan(*node.Node) error
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Sorter

// Sorter is an object that can sort a node.Node graph. The only requirements are
// that it is a stable sort.
type Sorter interface {
	Sort(*node.Node) error
}

// Collector is a type that can build a node.Node graph.
type Collector struct {
	scanner Scanner
	sorter  Sorter
	t       string
}

// New creates a new Collector.
func New(
	//registries []Registry,
	scanner Scanner,
	sorter Sorter,
	t string,
) *Collector {
	return &Collector{
		//registries: registries,
		scanner: scanner,
		sorter:  sorter,
		t:       t,
	}
}

// Collect creates a node.Node graph. It should return the node.Node that
// represents the target with which this Collector has been configured.
func (c *Collector) Collect() (*node.Node, error) {
	start := node.New(c.t)

	if err := c.s.Scan(start); err != nil {
		return nil, errors.Wrap(err, "scan")
	}

	return start, nil
}
