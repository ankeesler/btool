// Package builder provides functionality for resolving a full node.Node graph.
package builder

import (
	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Currenter

// Currenter tells whether a node.Node needs to be resolved.
type Currenter interface {
	Current(*node.Node) (bool, error)
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Callback

// Callback is an interface that clients can use to be notified about resolutions
// when they happen. It also tells the clients whether or not the node was up
// to date.
type Callback interface {
	OnResolve(n *node.Node, current bool)
}

// Builder can Build() a full node.Node graph.
type Builder struct {
	dryRun bool
	c      Currenter
	cb     Callback
}

// New creates a new Builder.
func New(dryRun bool, c Currenter, cb Callback) *Builder {
	return &Builder{
		dryRun: dryRun,
		c:      c,
		cb:     cb,
	}
}

func (b *Builder) Build(n *node.Node) error {
	return b.build(n, make(map[*node.Node]bool))
}

func (b *Builder) build(n *node.Node, built map[*node.Node]bool) error {
	if built[n] {
		return nil
	}

	log.Debugf("building %s", n.Name)

	for _, dN := range n.Dependencies {
		log.Debugf("building dependency %s", dN.Name)
		if err := b.build(dN, built); err != nil {
			return errors.Wrap(err, "resolve "+dN.Name)
		}
	}

	current, err := b.c.Current(n)
	if err != nil {
		return errors.Wrap(err, "current")
	}

	b.cb.OnResolve(n, current)

	if n.Resolver != nil {
		if !current {
			log.Debugf("resolving %s", n.Name)
			if !b.dryRun {
				if err := n.Resolver.Resolve(n); err != nil {
					return errors.Wrap(err, "really resolve "+n.Name)
				}
			}
		}
	}

	built[n] = true

	return nil
}
