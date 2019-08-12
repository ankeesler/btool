package registry

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

// FSRegistry holds Index()/Node() data in memory to be read by a client.
type FSRegistry struct {
	index *Index
	files map[string][]*Node
}

func newFSRegistry(index *Index, files map[string][]*Node) *FSRegistry {
	return &FSRegistry{
		index: index,
		files: files,
	}
}

// Create creates an FSRegistry from a directory. It will read files from a
// directory into memory. It returns an error if the FSRegistry cannot be created.
func CreateFSRegistry(fs afero.Fs, dir string) (*FSRegistry, error) {
	i := newIndex()
	files := make(map[string][]*Node)
	if err := afero.Walk(
		fs,
		dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !strings.HasSuffix(path, "_btool.yml") {
				logrus.Debugf("skipping path %s", path)
				return nil
			}

			data, err := afero.ReadFile(fs, path)
			if err != nil {
				return errors.Wrap(err, "read file")
			}

			nodes := make([]*Node, 0)
			if err := yaml.Unmarshal(data, &nodes); err != nil {
				return errors.Wrap(err, "unmarshal")
			}

			pathRel, err := filepath.Rel(dir, path)
			if err != nil {
				return errors.Wrap(err, "rel")
			}

			logrus.Debugf("adding file %s: %s", pathRel, nodes)
			i.Files = append(i.Files, IndexFile{
				Path:   pathRel,
				SHA256: sha256String(data),
			})
			files[pathRel] = nodes

			return nil
		},
	); err != nil {
		return nil, errors.Wrap(err, "walk")
	}

	logrus.Debugf("creating FSRegistry with %s/%s", i, files)
	return newFSRegistry(i, files), nil
}

func (fsr *FSRegistry) Index() (*Index, error) {
	return fsr.index, nil
}

func (fsr *FSRegistry) Nodes(name string) ([]*Node, error) {
	return fsr.files[name], nil
}

func sha256String(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum([]byte{}))
}
