package registry

import (
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// ClientCreator can create a Client.
type ClientCreator interface {
	Create() (Client, error)
}

// Creator is a type that creates a new Collector.
type Creator struct {
	fs    afero.Fs
	cc    ClientCreator
	cache string
	gc    GaggleCollector
}

// NewCreator will create a new Creator.
func NewCreator(
	fs afero.Fs,
	cc ClientCreator,
	cache string,
	gc GaggleCollector,
) *Creator {
	return &Creator{
		fs:    fs,
		cc:    cc,
		cache: cache,
		gc:    gc,
	}
}

// Create will create a new Collector. Is this godoc good enough?
func (c *Creator) Create() (*Collector, error) {
	client, err := c.cc.Create()
	if err != nil {
		return nil, errors.Wrap(err, "create")
	}

	return New(c.fs, client, c.cache, c.gc), nil
}
