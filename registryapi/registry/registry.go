package registry

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

// Package registry provides btool registry functionality.

// Registry holds files in memory to be read by a client.
type Registry struct {
	files map[string][]byte
}

func newRegistry() *Registry {
	return &Registry{
		files: make(map[string][]byte),
	}
}

// Create creates a Registry from a directory. It returns an error if the
// Registry cannot be created.
func Create(fs afero.Fs, dir string) (*Registry, error) {
	r := newRegistry()

	if err := afero.Walk(
		fs,
		dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			data, err := afero.ReadFile(fs, path)
			if err != nil {
				return errors.Wrap(err, "read file")
			}

			pathRel, err := filepath.Rel(dir, path)
			if err != nil {
				return errors.Wrap(err, "rel")
			}

			logrus.Debugf("adding file %s", pathRel)
			r.setFile(pathRel, data)

			return nil
		},
	); err != nil {
		return nil, errors.Wrap(err, "walk")
	}

	i, err := r.makeIndex()
	if err != nil {
		return nil, errors.Wrap(err, "make index")
	}
	r.setFile("index.yml", i)

	return r, nil
}

// Get returns an io.Reader for the provided file. It will return nil if no such
// file exists.
//
// Get callers should not cache the returned io.Reader. The returned io.Reader
// will be a new io.Reader every time this function is called. This is for the
// purpose of io.Reader cursor positions.
func (r *Registry) Get(file string) io.Reader {
	b, ok := r.files[file]
	if !ok {
		return nil
	} else {
		return bytes.NewBuffer(b)
	}
}

func (r *Registry) setFile(file string, data []byte) {
	r.files[file] = data
}

func (r *Registry) makeIndex() ([]byte, error) {
	i := Index{
		Files: make([]IndexFile, 0),
	}

	for file, data := range r.files {
		h := sha256.New()
		h.Write(data) // never returns an error, per doc
		i.Files = append(i.Files, IndexFile{
			Path:   file,
			SHA256: hex.EncodeToString(h.Sum([]byte{})),
		})
	}

	data, err := yaml.Marshal(&i)
	if err != nil {
		return nil, errors.Wrap(err, "marshal")
	}

	data = append([]byte("---\n"), data...)
	return data, nil
}
