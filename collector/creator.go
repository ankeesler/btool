package collector

import (
	"github.com/ankeesler/btool/log"
	"github.com/pkg/errors"
)

// CollectiniCreator is an object that can create a Collectini.
type CollectiniCreator interface {
	Create() (Collectini, error)
}

// Create can create a Collector.
type Creator struct {
	ctx    *Ctx
	cinics []CollectiniCreator
}

// NewCreator creates a new Creator.
func NewCreator(ctx *Ctx, cinics []CollectiniCreator) *Creator {
	return &Creator{
		ctx:    ctx,
		cinics: cinics,
	}
}

// Create will create a new Collector, injecting all of the Collectini's created
// via the CollectiniCreator's.
func (c *Creator) Create() (*Collector, error) {
	cinis := make([]Collectini, len(c.cinics))
	for i, cinic := range c.cinics {
		cini, err := cinic.Create()
		if err != nil {
			return nil, errors.Wrap(err, "create")
		}
		log.Debugf("created collectini #%d: %s", i, cini)
		cinis[i] = cini
	}
	return New(c.ctx, cinis...), nil
}