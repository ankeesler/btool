package handlers

import (
	"os"
	"path/filepath"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/pipeline"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type fs struct {
	effess afero.Fs
}

// NewFile creates pipeline.Handler that walks a file tree from a root and
// collects .c/.cc and .h files.
func NewFS(effess afero.Fs) pipeline.Handler {
	return &fs{
		effess: effess,
	}
}

func (fs *fs) Handle(ctx *pipeline.Ctx) {
	root := ctx.KV[pipeline.CtxRoot]

	logrus.Info("scanning from root " + root)

	if err := afero.Walk(
		fs.effess,
		root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return errors.Wrap(err, "walk")
			}

			rootRelPath, err := filepath.Rel(root, path)
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
			ctx.Nodes = append(ctx.Nodes, node.New(rootRelPath))

			return nil
		},
	); err != nil {
		ctx.Err = errors.Wrap(err, "walk")
	}
}

func (fs *fs) Name() string { return "fs" }
