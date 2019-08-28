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

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . GaggleCollector

// GaggleCollector is a type that can build a node.Node graph via a
// registrypkg.Gaggle. It is provided a root to indicate where the node.Node
// graph members should be located.
type GaggleCollector interface {
	Collect(ctx *collector.Ctx, g *registrypkg.Gaggle, root string) error
}

// Collector is a type that can build a node.Node graph via a btool registry.
type Collector struct {
	fs    afero.Fs
	r     registrypkg.Registry
	cache string
	gc    GaggleCollector
}

// New creates a new Collector.
func New(
	fs afero.Fs,
	r registrypkg.Registry,
	cache string,
	gc GaggleCollector,
) *Collector {
	return &Collector{
		fs:    fs,
		r:     r,
		cache: cache,
		gc:    gc,
	}
}

func (c *Collector) Collect(ctx *collector.Ctx, n *node.Node) error {
	i, err := c.r.Index()
	if err != nil {
		return errors.Wrap(err, "index")
	}

	for _, file := range i.Files {
		gaggleFile := filepath.Join(c.cache, file.SHA256+".yml")
		gaggle := new(registrypkg.Gaggle)
		log.Debugf("considering %s", gaggleFile)
		if exists, err := afero.Exists(c.fs, gaggleFile); err != nil {
			return errors.Wrap(err, "exists")
		} else if !exists {
			log.Debugf("does not exist")

			gaggle, err = c.r.Gaggle(file.Path)
			if err != nil {
				return errors.Wrap(err, "gaggle")
			} else if gaggle == nil {
				return errors.New("unknown gaggle at path: " + file.Path)
			}

			gaggleData, err := yaml.Marshal(&gaggle)
			if err != nil {
				return errors.Wrap(err, "marshal")
			}

			if err := c.fs.MkdirAll(filepath.Dir(gaggleFile), 0755); err != nil {
				return errors.Wrap(err, "mkdir all")
			}

			if err := afero.WriteFile(c.fs, gaggleFile, gaggleData, 0644); err != nil {
				return errors.Wrap(err, "write file")
			}
		} else {
			data, err := afero.ReadFile(c.fs, gaggleFile)
			if err != nil {
				return errors.Wrap(err, "read file")
			}

			if err := yaml.Unmarshal(data, &gaggle); err != nil {
				return errors.Wrap(err, "unmarshal")
			}
		}

		root := filepath.Join(c.cache, file.SHA256)
		if err := c.gc.Collect(ctx, gaggle, root); err != nil {
			return errors.Wrap(err, "collect")
		}
	}

	return nil
}
