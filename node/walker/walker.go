// Package walker provides a node.Handler that walks a file tree and creates nodes.
package walker

import (
	"os"
	"path/filepath"

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

func (w *Walker) Handle(cfg *node.Config, nodes []*node.Node) ([]*node.Node, error) {
	logrus.Info("scanning from root " + cfg.Root)

	if err := afero.Walk(
		w.fs,
		cfg.Root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return errors.Wrap(err, "walk")
			}

			rootRelPath, err := filepath.Rel(cfg.Root, path)
			if err != nil {
				return errors.Wrap(err, "rel")
			}

			if info.IsDir() {
				logrus.Debugf("skipping directory %s", rootRelPath)
				return nil
			} else {
				logrus.Debugf("looking at file %s", rootRelPath)
			}

			ext := filepath.Ext(rootRelPath)
			if ext != ".c" && ext != ".cc" && ext != ".h" {
				return nil
			}

			logrus.Debugf("adding node %s", rootRelPath)
			nodes = append(nodes, node.New(rootRelPath))

			return nil
		},
	); err != nil {
		return nil, errors.Wrap(err, "walk")
	}

	return nodes, nil
}
