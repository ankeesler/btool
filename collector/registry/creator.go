package registry

import (
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type RegistryCreator interface {
	Create() (Registry, error)
}

type Creator struct {
	fs    afero.Fs
	rc    RegistryCreator
	cache string
	g     Gaggler
}

func NewCreator(
	fs afero.Fs,
	rc RegistryCreator,
	cache string,
	g Gaggler,
) *Creator {
	return &Creator{
		fs:    fs,
		rc:    rc,
		cache: cache,
		g:     g,
	}
}

func (c *Creator) Create() (*Registry, error) {
	r, err := c.rc.Create()
	if err != nil {
		return errors.Wrap(err, "create")
	}

	return New(c.fs, r, c.ccache, c.g), nil
}
