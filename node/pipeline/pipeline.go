// Package pipeline provides an abtraction of a collection mechanism for
// node.Node's. A caller can create a Pipeline and add a bunch of Handler's
// that interact with a Ctx to collect node.Node's.
package pipeline

import (
	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Ctx

// Ctx is an interface with which a Handler interacts to collect node.Node's.
type Ctx interface {
	Add(*node.Node)
	Find(name string) *node.Node
	All() []*node.Node
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Handler

// Handler is an interface that describes some operator on a Pipeline.
//
// Upon being called, a Handler should add node.Node's to the provided Ctx, or
// return an error if it runs into trouble.
type Handler interface {
	Handle(Ctx) error
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Callback

// Callback is a notification mechanism for when node.Node's are added to a
// Ctx when running as part of a Pipeline.
type Callback interface {
	OnAdd(*node.Node)
}

// Pipeline is an abstraction of a list of operators on a mutable list of
// node.Node's.
type Pipeline struct {
	h  Handler
	cb Callback
}

// New creates a new Pipeline with a Handler and a Callback.
func New(h Handler, cb Callback) *Pipeline {
	return &Pipeline{
		h:  h,
		cb: cb,
	}
}

// Run kicks off the pipeline. It will return an error if any of the Handler's
// fail in their operation. It exits as soon as any Handler fails. If all
// Handler's succeed, then this function will return the Ctx on which the
// Handler's have been operating.
func (p *Pipeline) Run() (Ctx, error) {
	ctx := newCtx(p.cb)
	if err := p.h.Handle(ctx); err != nil {
		return nil, errors.Wrap(err, "handle")
	}
	return ctx, nil
}
