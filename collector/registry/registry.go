// Package registry provides functionality to build a node.Node graph using a
// btool registry.
package registry

import (
	"path/filepath"

	"github.com/ankeesler/btool/collector"
	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
	registrypkg "github.com/ankeesler/btool/registry"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Client

// Client is an object that can retrieve registrypkg.Gaggle's from a btool
// registry.
type Client interface {
	// Index should return the registrypkg.Index associated with this particular
	// Registry. If any error occurs, an error should be returned.
	Index() (*registrypkg.Index, error)
	// Gaggle should return the registrypkg.Gaggle associated with the provided
	// registrypkg.IndexFile.Path. If any error occurs, an error should be returned.
	// If no registrypkg.Gaggle exists for the provided string, then nil, nil should
	// be returned.
	Gaggle(string) (*registrypkg.Gaggle, error)
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . GaggleCollector

type GaggleCollector interface {
	Collect(ctx *collector.Ctx, g *registrypkg.Gaggle, root string) error
}

type Collector struct {
	fs    afero.Fs
	c     Client
	cache string
	gc    GaggleCollector
}

func New(
	fs afero.Fs,
	c Client,
	cache string,
	gc GaggleCollector,
) *Collector {
	return &Collector{
		fs:    fs,
		c:     c,
		cache: cache,
		gc:    gc,
	}
}

func (r *Collector) Collect(ctx *collector.Ctx, n *node.Node) error {
	i, err := r.c.Index()
	if err != nil {
		return errors.Wrap(err, "index")
	}

	for _, file := range i.Files {
		gaggleFile := filepath.Join(r.cache, file.SHA256+".yml")
		gaggle := new(registrypkg.Gaggle)
		log.Debugf("considering %s", gaggleFile)
		if exists, err := afero.Exists(r.fs, gaggleFile); err != nil {
			return errors.Wrap(err, "exists")
		} else if !exists {
			log.Debugf("does not exist")

			gaggle, err = r.c.Gaggle(file.Path)
			if err != nil {
				return errors.Wrap(err, "gaggle")
			} else if gaggle == nil {
				return errors.New("unknown gaggle at path: " + file.Path)
			}

			gaggleData, err := yaml.Marshal(&gaggle)
			if err != nil {
				return errors.Wrap(err, "marshal")
			}

			if err := r.fs.MkdirAll(filepath.Dir(gaggleFile), 0755); err != nil {
				return errors.Wrap(err, "mkdir all")
			}

			if err := afero.WriteFile(r.fs, gaggleFile, gaggleData, 0644); err != nil {
				return errors.Wrap(err, "write file")
			}
		} else {
			data, err := afero.ReadFile(r.fs, gaggleFile)
			if err != nil {
				return errors.Wrap(err, "read file")
			}

			if err := yaml.Unmarshal(data, &gaggle); err != nil {
				return errors.Wrap(err, "unmarshal")
			}
		}

		root := filepath.Join(r.cache, file.SHA256)
		if err := r.gc.Collect(ctx, gaggle, root); err != nil {
			return errors.Wrap(err, "collect")
		}
	}

	return nil
}
