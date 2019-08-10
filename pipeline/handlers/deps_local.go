package handlers

import (
	"fmt"
	"path/filepath"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/pipeline"
	"github.com/ankeesler/btool/pipeline/handlers/includes"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type depsLocal struct {
	fs afero.Fs
}

// NewDepsLocal returns a pipeline.Handler that collects the local dependencies
// for each node.Node in the node list.
func NewDepsLocal(fs afero.Fs) pipeline.Handler {
	return &depsLocal{
		fs: fs,
	}
}

func (dl *depsLocal) Handle(ctx *pipeline.Ctx) error {
	nodeMap := make(map[string]*node.Node)
	for _, n := range ctx.Nodes {
		nodeMap[n.Name] = n
	}

	root := ctx.KV[pipeline.CtxRoot]
	for _, n := range ctx.Nodes {
		logrus.Debugf("deps_local: handling node %s", n)
		if err := dl.handleNode(n, nodeMap, root); err != nil {
			return errors.Wrap(err, fmt.Sprintf("handle node %s", n.Name))
		}
	}

	return nil
}

func (dl *depsLocal) Name() string { return "local deps" }

func (dl *depsLocal) handleNode(
	n *node.Node,
	nodeMap map[string]*node.Node,
	root string,
) error {
	path := filepath.Join(root, n.Name)
	data, err := afero.ReadFile(dl.fs, path)
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
		includeName, err := dl.resolveInclude(include, filepath.Dir(path), root)
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
func (dl *depsLocal) resolveInclude(include, dir, root string) (string, error) {
	rootRelJoin := filepath.Join(root, include)
	if exists, err := afero.Exists(dl.fs, rootRelJoin); err != nil {
		return "", errors.Wrap(err, "exists")
	} else if exists {
		return include, nil
	}

	dirRelJoin := filepath.Join(dir, include)
	if exists, err := afero.Exists(dl.fs, dirRelJoin); err != nil {
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
