// Package app provides the btool application that can be used to perform
// btool domain work.
package app

import (
	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Collector

// Collector creates a node.Node graph.
type Collector interface {
	Collect(*node.Node) error
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . CollectorCreator

// CollectorCreator creates a Collector.
type CollectorCreator interface {
	Create() (Collector, error)
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

// App is a type that does the domain work of a btool invocation.
type App struct {
	cc      CollectorCreator
	cleaner Cleaner
	builder Builder
}

// New returns a new App struct.
func New(cc CollectorCreator, cleaner Cleaner, builder Builder) *App {
	return &App{
		cc:      cc,
		cleaner: cleaner,
		builder: builder,
	}
}

// Run runs a btool build/clean.
func (a *App) Run(n *node.Node, clean, dryRun bool) error {
	c, err := a.cc.Create()
	if err != nil {
		return errors.Wrap(err, "create")
	}

	if err := c.Collect(n); err != nil {
		return errors.Wrap(err, "collect")
	}

	log.Debugf("graph: %s", node.String(n))

	if clean {
		if err := a.cleaner.Clean(n); err != nil {
			return errors.Wrap(err, "clean")
		}
	} else {
		if err := a.builder.Build(n); err != nil {
			return errors.Wrap(err, "build")
		}
	}

	return nil
}
