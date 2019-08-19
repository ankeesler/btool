package pipeline

import (
	"fmt"

	"github.com/pkg/errors"
)

// MultiHandler is a Handler that calls Handle() on a slice of Handler's.
type MultiHandler struct {
	hs []Handler
}

// NewMultiHandler creates a MultiHandler with no Handler's.
func NewMultiHandler() *MultiHandler {
	return &MultiHandler{}
}

// Add adds one or more Handler's to the MultiHandler. It returns the
// MultiHandler so that calls can be strung together.
//   mh := NewMultiHandler().Add(h0, h1).Add(h2)
func (mh *MultiHandler) Add(h ...Handler) *MultiHandler {
	mh.hs = append(mh.hs, h...)
	return mh
}

// Handle calls Handle() on each of the MultiHandler's Handler's.
func (mh *MultiHandler) Handle(ctx Ctx) error {
	for _, h := range mh.hs {
		if err := h.Handle(ctx); err != nil {
			return errors.Wrap(err, fmt.Sprintf("handle: %s", h))
		}
	}
	return nil
}
