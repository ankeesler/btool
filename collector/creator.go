package collector

import "github.com/pkg/errors"

type CollectiniCreator interface {
	Create() (Collectini, error)
}

type Creator struct {
	ctx    *Ctx
	cinics []CollectiniCreator
}

func NewCreator(ctx *Ctx, cinics []CollectiniCreator) *Creator {
	return &Creator{
		ctx:    ctx,
		cinics: cinics,
	}
}

func (c *Creator) Create() (*Collector, error) {
	cinis := make([]Collectini, len(c.cinics))
	for i := range c.cinics {
		cini, err := c.cinics[i].Create()
		if err != nil {
			return nil, errors.Wrap(err, "create")
		}
		cinis = append(cinis, cini)
	}
	return New(c.ctx, cinis...), nil
}
