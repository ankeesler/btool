package registry

import (
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type ClientCreator interface {
	Create() (Client, error)
}

type Creator struct {
	fs    afero.Fs
	cc    ClientCreator
	cache string
	gc    GaggleCollector
}

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

func (c *Creator) Create() (*Collector, error) {
	client, err := c.cc.Create()
	if err != nil {
		return nil, errors.Wrap(err, "create")
	}

	return New(c.fs, client, c.cache, c.gc), nil
}
