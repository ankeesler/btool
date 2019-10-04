package registry

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"

	"github.com/ankeesler/btool/log"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

// FSRegistry holds Index/Gaggle data in memory to be read by a client.
type FSRegistry struct {
	index *Index
	files map[string]*Gaggle
}

func newFSRegistry(index *Index, files map[string]*Gaggle) *FSRegistry {
	return &FSRegistry{
		index: index,
		files: files,
	}
}

// CreateFSRegistry creates an FSRegistry from a directory. It will read files from a
// directory into memory. It returns an error if the FSRegistry cannot be created.
func CreateFSRegistry(fs afero.Fs, dir, name string) (*FSRegistry, error) {
	i := newIndex(name)
	files := make(map[string]*Gaggle)
	if err := afero.Walk(
		fs,
		dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if path == "index.yml" || !strings.HasSuffix(path, ".yml") {
				log.Debugf("skipping path %s", path)
				return nil
			}

			data, err := afero.ReadFile(fs, path)
			if err != nil {
				return errors.Wrap(err, "read file")
			}

			gaggle := new(Gaggle)
			if err := yaml.Unmarshal(data, gaggle); err != nil {
				return errors.Wrap(err, "unmarshal")
			}

			pathRel, err := filepath.Rel(dir, path)
			if err != nil {
				return errors.Wrap(err, "rel")
			}

			log.Debugf("adding file %s: %s", pathRel, gaggle)
			i.Files = append(i.Files, IndexFile{
				Path:   pathRel,
				SHA256: sha256String(data),
			})
			files[pathRel] = gaggle

			return nil
		},
	); err != nil {
		return nil, errors.Wrap(err, "walk")
	}

	log.Debugf("creating FSRegistry with %s/%v", i, files)
	return newFSRegistry(i, files), nil
}

func (fsr *FSRegistry) Index() (*Index, error) {
	return fsr.index, nil
}

func (fsr *FSRegistry) Gaggle(name string) (*Gaggle, error) {
	return fsr.files[name], nil
}

func sha256String(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum([]byte{}))
}
