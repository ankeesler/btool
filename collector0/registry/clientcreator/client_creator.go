// Package clientcreator provides functionality for creating a registry.Client.
package clientcreator

import (
	"net/http"
	"net/url"

	"github.com/ankeesler/btool/collector/registry"
	"github.com/ankeesler/btool/log"
	registrypkg "github.com/ankeesler/btool/registry"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// Creator is a type that can create a registry.Client
type Creator struct {
	fs  afero.Fs
	url string
}

// New returns a new Creator.
func New(fs afero.Fs, url string) *Creator {
	return &Creator{
		fs:  fs,
		url: url,
	}
}

// Create will create a registry.Client.
func (c *Creator) Create() (registry.Client, error) {
	url, err := url.Parse(c.url)
	if err != nil {
		return nil, errors.Wrap(err, "url parse")
	}

	var client registry.Client
	switch url.Scheme {
	case "http", "https":
		client = registrypkg.NewHTTPRegistry(c.url, &http.Client{})
		log.Debugf("creating http registry from %s", url.Host)
	case "file":
		client, err = registrypkg.CreateFSRegistry(c.fs, url.Path, c.url)
		log.Debugf("creating fs registry from %s", url.Path)
	default:
		client, err = registrypkg.CreateFSRegistry(c.fs, c.url, c.url)
		log.Debugf("creating fs registry from %s", c.url)
	}

	return client, err
}
