package handlers

import (
	"fmt"
	"path/filepath"

	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/pipeline"
	"github.com/ankeesler/btool/node/pipeline/handlers/includes"
	"github.com/ankeesler/btool/node/pipeline/handlers/walk"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type fs struct {
	effess afero.Fs
	root   string
}

// NewFile creates pipeline.Handler that walks a file tree from a root and
// collects .c/.cc and .h files.
func NewFS(effess afero.Fs, root string) pipeline.Handler {
	return &fs{
		effess: effess,
		root:   root,
	}
}

func (fs *fs) Handle(ctx *pipeline.Ctx) error {
	log.Debugf("scanning from root %s", fs.root)

	nodeMap := make(map[string]*node.Node)
	if err := walk.Walk(
		fs.effess,
		fs.root,
		[]string{".c", ".cc", ".h"},
		func(file string) error {
			log.Debugf("adding node %s", file)
			n := node.New(file)
			ctx.Nodes = append(ctx.Nodes, n)
			nodeMap[n.Name] = n

			return nil
		},
	); err != nil {
		return errors.Wrap(err, "walk")
	}

	for _, n := range nodeMap {
		log.Debugf("deps_local: handling node %s", n)
		if err := fs.handleNode(n, nodeMap, fs.root); err != nil {
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
	path := n.Name
	data, err := afero.ReadFile(fs.effess, path)
	if err != nil {
		return errors.Wrap(err, "read file "+path)
	}
	log.Debugf("read file %s", path)

	includes, err := includes.Parse(data)
	if err != nil {
		return errors.Wrap(err, "parse includes")
	}
	log.Debugf("includes = %s", includes)

	for _, include := range includes {
		includeName, err := fs.resolveInclude(include, filepath.Dir(path), root)
		if err != nil {
			return errors.Wrap(err, "resolve include path "+include)
		}

		includeN, ok := nodeMap[includeName]
		if !ok {
			return fmt.Errorf("unknown node for include name %s", includeName)
		}

		log.Debugf("adding dependency %s -> %s", n.Name, includeN.Name)
		n.Dependencies = append(n.Dependencies, includeN)
	}

	return nil
}

func (fs *fs) resolveInclude(include, dir, root string) (string, error) {
	rootRelJoin := filepath.Join(root, include)
	if exists, err := afero.Exists(fs.effess, rootRelJoin); err != nil {
		return "", errors.Wrap(err, "exists")
	} else if exists {
		return rootRelJoin, nil
	}

	dirRelJoin := filepath.Join(dir, include)
	if exists, err := afero.Exists(fs.effess, dirRelJoin); err != nil {
		return "", errors.Wrap(err, "exists")
	} else if exists {
		if err != nil {
			return "", errors.New("rel")
		} else {
			return dirRelJoin, nil
		}
	}

	return "", errors.New("cannot resolve include: " + include)
}
