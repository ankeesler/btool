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
	Collect() error
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

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Lister

// Lister prints out the members of the node.Node graph.
type Lister interface {
	List(*node.Node) error
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Builder

// Builder brings the node.Node graph into existence.
type Builder interface {
	Build(*node.Node) error
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Runner

// Runner runs a node.Node.
type Runner interface {
	Run(*node.Node) error
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Watcher

// Watcher listens for changes in the node.Node graph. The call will block until
// a change is made to one of the provided node.Node's on disk. The call will
// return an error if there was a problem with the call.
type Watcher interface {
	Watch(*node.Node) error
}

// App is a type that does the domain work of a btool invocation.
type App struct {
	cc      CollectorCreator
	cleaner Cleaner
	lister  Lister
	builder Builder
	runner  Runner
	watcher Watcher
}

// New returns a new App struct.
func New(
	cc CollectorCreator,
	cleaner Cleaner,
	lister Lister,
	builder Builder,
	runner Runner,
	watcher Watcher,
) *App {
	return &App{
		cc:      cc,
		cleaner: cleaner,
		lister:  lister,
		builder: builder,
		runner:  runner,
		watcher: watcher,
	}
}

// Run runs a btool build/clean.
func (a *App) Run(n *node.Node, clean, list, run, watch bool) error {
	c, err := a.cc.Create()
	if err != nil {
		return errors.Wrap(err, "create")
	}

	if err := c.Collect(); err != nil {
		return errors.Wrap(err, "collect")
	}

	log.Debugf("graph: %s", node.String(n))

	if clean {
		if err := a.cleaner.Clean(n); err != nil {
			return errors.Wrap(err, "clean")
		}
	} else if list {
		if err := a.lister.List(n); err != nil {
			return errors.Wrap(err, "list")
		}
	} else {
		for {
			err := a.builder.Build(n)
			if err != nil {
				if !watch {
					return errors.Wrap(err, "build")
				} else {
					log.Errorf("build: %s", err)
				}
			}

			if err == nil && run {
				if err := a.runner.Run(n); err != nil {
					if !watch {
						return errors.Wrap(err, "run")
					} else {
						log.Errorf("run: %s", err)
					}
				}
			}

			if watch {
				if err := a.watcher.Watch(n); err != nil {
					return errors.Wrap(err, "watch")
				}
			} else {
				break
			}
		}
	}

	return nil
}
