package deps

import (
	"fmt"
	"path/filepath"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/deps/includes"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// Resolves dependencies local to project.
type Local struct {
	fs afero.Fs
}

func NewLocal(fs afero.Fs) *Local {
	return &Local{
		fs: fs,
	}
}

func (l *Local) Handle(cfg *node.Config, nodes []*node.Node) ([]*node.Node, error) {
	nodeMap := make(map[string]*node.Node)
	for _, node := range nodes {
		nodeMap[node.Name] = node
	}

	for _, node := range nodes {
		logrus.Debugf("local_deps: handling node %s", node)
		if err := l.handleNode(cfg, node, nodeMap); err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("handle node %s", node.Name))
		}
	}
	return nodes, nil
}

func (l *Local) handleNode(
	cfg *node.Config,
	node *node.Node,
	nodeMap map[string]*node.Node,
) error {
	path := filepath.Join(cfg.Root, node.Name)
	data, err := afero.ReadFile(l.fs, path)
	if err != nil {
		return errors.Wrap(err, "read file "+path)
	}
	logrus.Debugf("read file %s", path)

	includes, err := includes.Parse(data)
	if err != nil {
		return errors.Wrap(err, "parse includes")
	}

	for _, include := range includes {
		includeName, err := l.resolveInclude(cfg, include, filepath.Dir(path))
		if err != nil {
			return errors.Wrap(err, "resolve include path "+include)
		}

		includeNode, ok := nodeMap[includeName]
		if !ok {
			return fmt.Errorf("unknown node for include name %s", includeName)
		}

		node.Dependencies = append(node.Dependencies, includeNode)
	}

	return nil
}

// Return an include path relative to the root!
func (l *Local) resolveInclude(cfg *node.Config, include, dir string) (string, error) {
	rootRelJoin := filepath.Join(cfg.Root, include)
	if exists, err := afero.Exists(l.fs, rootRelJoin); err != nil {
		return "", errors.Wrap(err, "exists")
	} else if exists {
		return include, nil
	}

	dirRelJoin := filepath.Join(dir, include)
	if exists, err := afero.Exists(l.fs, dirRelJoin); err != nil {
		return "", errors.Wrap(err, "exists")
	} else if exists {
		rootRelJoin, err := filepath.Rel(cfg.Root, dirRelJoin)
		if err != nil {
			return "", errors.New("rel")
		} else {
			return rootRelJoin, nil
		}
	}

	return "", errors.New("cannot resolve include: " + include)
}
