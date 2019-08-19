package handlers

import (
	"fmt"
	"path/filepath"

	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/pipeline"
	"github.com/pkg/errors"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Collector

// Collector is a type that can make a list of files from a root directory that
// match some list of file extensions. It should follow symlinks.
type Collector interface {
	Collect(root string, exts []string) ([]string, error)
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Includeser

// Includeser is a type that can return a list of #include's from a given file.
type Includeser interface {
	Includes(path string) ([]string, error)
}

type fs struct {
	c    Collector
	i    Includeser
	root string
}

// NewFS creates a pipeline.Handler that walks a file tree from a root and
// collects .c/.cc and .h files.
func NewFS(
	c Collector,
	i Includeser,
	root string,
) pipeline.Handler {
	return &fs{
		c:    c,
		i:    i,
		root: root,
	}
}

func (fs *fs) Handle(ctx pipeline.Ctx) error {
	log.Debugf("scanning from root %s", fs.root)

	nodeMap := make(map[string]*node.Node)
	paths, err := fs.c.Collect(fs.root, []string{".c", ".cc", ".h"})
	if err != nil {
		return errors.Wrap(err, "collect")
	}

	for _, path := range paths {
		log.Debugf("adding node %s", path)
		n := node.New(path)
		ctx.Nodes = append(ctx.Nodes, n)
		nodeMap[path] = n
	}

	for _, n := range nodeMap {
		log.Debugf("deps_local: handling node %s", n)
		if err := fs.handleNode(n, nodeMap, fs.root); err != nil {
			return errors.Wrap(err, fmt.Sprintf("handle node %s", n.Name))
		}
	}

	return nil
}

func (fs *fs) String() string { return "fs" }

func (fs *fs) handleNode(
	n *node.Node,
	nodeMap map[string]*node.Node,
	root string,
) error {
	path := n.Name
	includes, err := fs.i.Includes(path)
	if err != nil {
		return errors.Wrap(err, "includes")
	}
	log.Debugf("includes = %s", includes)

	for _, include := range includes {
		includeN, ok := nodeMap[filepath.Join(root, include)]
		if !ok {
			includeN, ok = nodeMap[filepath.Join(filepath.Dir(path), include)]
		}
		if !ok {
			return fmt.Errorf("unknown node for include %s", include)
		}

		log.Debugf("adding dependency %s -> %s", n.Name, includeN.Name)
		n.Dependencies = append(n.Dependencies, includeN)
	}

	return nil
}
