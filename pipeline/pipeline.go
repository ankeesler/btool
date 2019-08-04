// Package pipeline provides an abstraction of a series of operators (i.e.,
// Handler's) on a list of node.Node's.
package pipeline

import (
	"fmt"

	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Handler

// Handler is an interface that describes some operator on a Pipeline. The
// Handler should perform work on the provided list of node.Node's and return the
// new list. It should return an error if there was a problem while performing
// its operation.
type Handler interface {
	Handle(*Ctx, []*node.Node) ([]*node.Node, error)

	// Name returns an identifying name for this Handler. This helps with
	// debugging.
	Name() string
}

// Pipeline is an abstraction of a list of operators on a mutable list of
// node.Node's.
type Pipeline struct {
	ctx      *Ctx
	handlers []Handler
}

// New creates a new Pipeline with a Ctx and a list of Handler's.
func New(ctx *Ctx, handlers ...Handler) *Pipeline {
	return &Pipeline{
		ctx:      ctx,
		handlers: handlers,
	}
}

// Run kicks off the pipeline. It will return an error if any of the Handler's
// fail in their operation. It exits as soon as any Handler fails.
func (p *Pipeline) Run() error {
	nodes := make([]*node.Node, 0)

	var err error
	for _, h := range p.handlers {
		logrus.Debugf("pipeline: running %s", h.Name())
		nodes, err = h.Handle(p.ctx, nodes)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("handle (%s)", h.Name()))
		}
	}

	return nil
}
