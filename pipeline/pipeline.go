// Package pipeline provides an abstraction of a series of operators (i.e.,
// Handler's) on a list of node.Node's.
package pipeline

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Handler

// Handler is an interface that describes some operator on a Pipeline.
//
// Upon being called, a Handler should:
//   - update the node.Node list field of the provided Ctx in order to propagate
//     their updates to the node collection
//   - set the err field of the provided Ctx if it runs into an error
//   - get or set any keys on the Ctx to provide information to other Handler's
//     in the Pipeline
type Handler interface {
	Handle(*Ctx)

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
	for _, h := range p.handlers {
		logrus.Debugf("pipeline: running %s", h.Name())
		h.Handle(p.ctx)
		if p.ctx.Err != nil {
			return errors.Wrap(p.ctx.Err, fmt.Sprintf("handle (%s)", h.Name()))
		}
	}

	return nil
}
