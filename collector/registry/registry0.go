package registry

import (
	"path/filepath"

	"github.com/ankeesler/btool/collector"
	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/registry"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Registry

// Registry is an object that can retrieve registrypkg.Gaggle's.
type Registry interface {
	// Index should return the registrypkg.Index associated with this particular
	// Registry. If any error occurs, an error should be returned.
	Index() (*registrypkg.Index, error)
	// Gaggle should return the registrypkg.Gaggle associated with the provided
	// registrypkg.IndexFile.Path. If any error occurs, an error should be returned.
	// If no registrypkg.Gaggle exists for the provided string, then nil, nil should
	// be returned.
	Gaggle(string) (*registrypkg.Gaggle, error)
}

type Gaggler interface {
	Collect(ctx *collector.Ctx, g *registrypkg.Gaggle)
}

type Registry struct {
	fs    afero.Fs
	r     Registry
	cache string
	g     Gaggler
}

func New(
	fs afero.Fs,
	r Registry,
	cache string,
	g Gaggler,
) *Registry {
	return &Registry{
		fs:    fs,
		r:     r,
		cache: cache,
		g:     g,
	}
}

func (r *Registry) Collect(ctx *collector.Ctx, n *node.Node) error {
	i, err := r.r.Index()
	if err != nil {
		return errors.Wrap(err, "index")
	}

	for i, file := range i.Files {
		gaggleFile := filepath.Join(f.cache, file.SHA256+".yml")
		gaggle := new(registry.Gaggle)
		log.Debugf("considering %s", gaggleFile)
		if exists, err := afero.Exists(f.fs, gaggleFile); err != nil {
			return errors.Wrap(err, "exists")
		} else if !exists {
			log.Debugf("does not exist")

			gaggle, err = f.r.Gaggle(file.Path)
			if err != nil {
				return errors.Wrap(err, "gaggle")
			} else if gaggle == nil {
				return errors.New("unknown gaggle at path: " + file.Path)
			}

			gaggleData, err := yaml.Marshal(&gaggle)
			if err != nil {
				return errors.Wrap(err, "marshal")
			}

			if err := f.fs.MkdirAll(filepath.Dir(gaggleFile), 0755); err != nil {
				return errors.Wrap(err, "mkdir all")
			}

			if err := afero.WriteFile(f.fs, gaggleFile, gaggleData, 0644); err != nil {
				return errors.Wrap(err, "write file")
			}
		} else {
			data, err := afero.ReadFile(f.fs, gaggleFile)
			if err != nil {
				return errors.Wrap(err, "read file")
			}

			if err := yaml.Unmarshal(data, &gaggle); err != nil {
				return errors.Wrap(err, "unmarshal")
			}
		}

		if err := f.g.Collect(ctx, gaggle); err != nil {
			return errors.Wrap(err, "collect")
		}
	}

	return nil
}
