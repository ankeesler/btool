// Package pipeline provides an abstraction of a series of operators (i.e.,
// Handler's) on a list of node.Node's.
package pipeline

import (
	"fmt"

	"github.com/ankeesler/btool/log"
	"github.com/pkg/errors"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Handler

// Handler is an interface that describes some operator on a Pipeline.
//
// Upon being called, a Handler should:
//   - update the node.Node list field of the provided Ctx in order to propagate
//     their updates to the node collection
//   - get or set any keys on the Ctx to provide information to other Handler's
//     in the Pipeline
//   - return an error if something goes wrong, skrrrt
type Handler interface {
	Handle(*Ctx) error

	// TODO: change me to String()! And get rid of me!
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

// Handlers is a builder method that allows a caller to add a Handler to the
// Pipeline. It returns the Pipeline so that the calls can be strung together.
//   p := New(handlerA, handlerB).Handler(handlerC).Handler(handlerD)
func (p *Pipeline) Handlers(handlers ...Handler) *Pipeline {
	for _, h := range handlers {
		p.handlers = append(p.handlers, h)
	}
	return p
}

// Run kicks off the pipeline. It will return an error if any of the Handler's
// fail in their operation. It exits as soon as any Handler fails.
func (p *Pipeline) Run() error {
	for _, h := range p.handlers {
		log.Debugf("pipeline: running %s", h.Name())
		if err := h.Handle(p.ctx); err != nil {
			return errors.Wrap(err, fmt.Sprintf("handle (%s)", h.Name()))
		}
	}

	return nil
}
