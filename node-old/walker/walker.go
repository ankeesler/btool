// Package walker provides a node.Handler that walks a file tree and creates nodes.
package walker

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type Walker struct {
	fs afero.Fs
}

// New creates a new Walker with a root directory. The root directory will be the
// starting place for the walk.
func New(fs afero.Fs) *Walker {
	return &Walker{
		fs: fs,
	}
}

func (w *Walker) Handle(c *node.Config, nodes []*node.Node) ([]*node.Node, error) {
	logrus.Info("scanning from root " + c.Root)

	if err := afero.Walk(
		w.fs,
		c.Root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return errors.Wrap(err, "walk")
			}

			rootRelPath, err := filepath.Rel(c.Root, path)
			if err != nil {
				return errors.Wrap(err, "rel")
			}

			if info.IsDir() {
				logrus.Debugf("skipping directory %s", rootRelPath)
				return nil
			} else {
				logrus.Debugf("looking at file %s", rootRelPath)
			}

			nodes = append(nodes, &node.Node{
				Name:    rootRelPath,
				Sources: sources(rootRelPath),
				Headers: headers(rootRelPath),
			})

			return nil
		},
	); err != nil {
		return nil, errors.Wrap(err, "walk")
	}

	return nodes, nil
}

func sources(name string) []string {
	if strings.HasSuffix(name, ".c") || strings.HasSuffix(name, ".cc") {
		return []string{name}
	} else {
		return []string{}
	}
}

func headers(name string) []string {
	if strings.HasSuffix(name, ".h") {
		return []string{name}
	} else {
		return []string{}
	}
}