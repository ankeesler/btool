// Package btool provides fundamental entities that can be used to perform
// btool domain work.
//
// Clients that are lazy should call Run() with the desired Cfg.
//   ...
//   err := Run(&Cfg{
//     ...
//   })
//   ...
//
// Clients want to go above and beyong should call New() and then Run() on the
// returned Btool struct.
//   ...
//   b := New(...)
//   err := b.Run(...)
//   ...
//
package btool

import (
	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Collector

// Collector creates a node.Node graph.
type Collector interface {
	Collect(*node.Node) (*node.Node, error)
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Cleaner

// Cleaner removes the node.Node graph from disk.
type Cleaner interface {
	Clean(*node.Node) error
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Builder

// Builder brings the node.Node graph into existence.
type Builder interface {
	Build(*node.Node) error
}

// Btool is a type that does the domain work of a btool invocation.
type Btool struct {
	collector Collector
	cleaner   Cleaner
	builder   Builder
}

// New returns a new Btool struct.
func New(collector Collector, cleaner Cleaner, builder Builder) *Btool {
	return &Btool{
		collector: collector,
		cleaner:   cleaner,
		builder:   builder,
	}
}

// Run runs a btool build/clean.
func (b *Btool) Run(n *node.Node, clean, dryRun bool) error {
	var err error

	n, err = b.collector.Collect(n)
	if err != nil {
		return errors.Wrap(err, "collect")
	}

	log.Debugf("graph: %s", node.String(n))

	if clean {
		if err := b.cleaner.Clean(n); err != nil {
			return errors.Wrap(err, "clean")
		}
	} else {
		if err := b.builder.Build(n); err != nil {
			return errors.Wrap(err, "build")
		}
	}

	return nil
}
