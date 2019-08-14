package handlers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/pipeline"
	"github.com/ankeesler/btool/node/pipeline/handlers/includes"
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

func (fs *fs) Handle(ctx *pipeline.Ctx) error {
	root := ctx.KV[pipeline.CtxRoot]

	logrus.Debugf("scanning from root %s", root)

	nodeMap := make(map[string]*node.Node)
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
			n := node.New(rootRelPath)
			ctx.Nodes = append(ctx.Nodes, n)
			nodeMap[n.Name] = n

			return nil
		},
	); err != nil {
		return errors.Wrap(err, "walk")
	}

	for _, n := range nodeMap {
		logrus.Debugf("deps_local: handling node %s", n)
		if err := fs.handleNode(n, nodeMap, root); err != nil {
			return errors.Wrap(err, fmt.Sprintf("handle node %s", n.Name))
		}
	}

	return nil
}

func (fs *fs) Name() string { return "fs" }

func (fs *fs) handleNode(
	n *node.Node,
	nodeMap map[string]*node.Node,
	root string,
) error {
	path := filepath.Join(root, n.Name)
	data, err := afero.ReadFile(fs.effess, path)
	if err != nil {
		return errors.Wrap(err, "read file "+path)
	}
	logrus.Debugf("read file %s", path)

	includes, err := includes.Parse(data)
	if err != nil {
		return errors.Wrap(err, "parse includes")
	}
	logrus.Debugf("includes = %s", includes)

	for _, include := range includes {
		includeName, err := fs.resolveInclude(include, filepath.Dir(path), root)
		if err != nil {
			return errors.Wrap(err, "resolve include path "+include)
		}

		includeN, ok := nodeMap[includeName]
		if !ok {
			return fmt.Errorf("unknown node for include name %s", includeName)
		}

		logrus.Debugf("adding dependency %s -> %s", n.Name, includeN.Name)
		n.Dependencies = append(n.Dependencies, includeN)
	}

	return nil
}

// Return an include path relative to the root!
func (fs *fs) resolveInclude(include, dir, root string) (string, error) {
	rootRelJoin := filepath.Join(root, include)
	if exists, err := afero.Exists(fs.effess, rootRelJoin); err != nil {
		return "", errors.Wrap(err, "exists")
	} else if exists {
		return include, nil
	}

	dirRelJoin := filepath.Join(dir, include)
	if exists, err := afero.Exists(fs.effess, dirRelJoin); err != nil {
		return "", errors.Wrap(err, "exists")
	} else if exists {
		rootRelJoin, err := filepath.Rel(root, dirRelJoin)
		if err != nil {
			return "", errors.New("rel")
		} else {
			return rootRelJoin, nil
		}
	}

	return "", errors.New("cannot resolve include: " + include)
}