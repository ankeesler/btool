package gaggler

import (
	"path/filepath"

	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/registry"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Registry

// Registry is an object that can retrieve registry.Gaggle's.
type Registry interface {
	// Index should return the registry.Index associated with this particular
	// Registry. If any error occurs, an error should be returned.
	Index() (*registry.Index, error)
	// Gaggle should return the registry.Gaggle associated with the provided
	// registry.IndexFile.Path. If any error occurs, an error should be returned.
	// If no registry.Gaggle exists for the provided string, then nil, nil should
	// be returned.
	Gaggle(string) (*registry.Gaggle, error)
}

// Factory is a type that create Gaggler's.
//
// It caches Gaggle()'s by their SHA256 sum.
// It always retrieves the Index() from the Registry, but it may not have to
// retrieve any of the the Gaggle()'s associated with that Registry.
type Factory struct {
	fs    afero.Fs
	r     Registry
	cache string

	gaggles      []*registry.Gaggle
	roots        []string
	gagglesIndex int
}

// NewFactory creates a new Factory.
func NewFactory(fs afero.Fs, r Registry, cache string) *Factory {
	return &Factory{
		fs:    fs,
		r:     r,
		cache: cache,

		gaggles: nil,
	}
}

// Next iterates through the Gaggle()'s that this Factory has gotten from a
// Registry and creates a Gaggler for each Gaggle(). It returns the sequence of
// Gaggler's. OOnce all of the Gaggler's have been iterated through, this
// function will return a nil Gaggler.
//
// It can be used like this.
//   f := NewFactory(...)
//   for {
//     g, err := f.Next()
//     if err != nil {
//       // handle err...
//     } else if g == nil {
//       break
//     } else {
//       // handle g...
//     }
//   }
func (f *Factory) Next() (*Gaggler, error) {
	if f.gaggles == nil {
		if err := f.initGaggles(); err != nil {
			return nil, errors.Wrap(err, "init gaggle num")
		}
	} else if f.gagglesIndex >= len(f.gaggles) {
		return nil, nil
	}

	g := f.gaggles[f.gagglesIndex]
	r := f.roots[f.gagglesIndex]
	f.gagglesIndex++
	return newGaggler(g, r), nil
}

func (f *Factory) initGaggles() error {
	i, err := f.r.Index()
	if err != nil {
		return errors.Wrap(err, "index")
	}

	f.gaggles = make([]*registry.Gaggle, len(i.Files))
	f.roots = make([]string, len(i.Files))

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

		f.gaggles[i] = gaggle
		f.roots[i] = filepath.Join(f.cache, file.SHA256)
	}

	return nil
}
