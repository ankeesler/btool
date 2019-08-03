package deps

import (
	"path/filepath"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/deps/includes"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Downloader

type Downloader interface {
	Download(afero.Fs, string, string, string) error
}

// Resolves dependencies outside of a project.
type Remote struct {
	fs    afero.Fs
	root  string
	cache string
	d     Downloader
}

func NewRemote(fs afero.Fs, root string, cache string, d Downloader) *Remote {
	return &Remote{
		fs:    fs,
		root:  root,
		cache: cache,
		d:     d,
	}
}

func (r *Remote) Handle(nodes []*node.Node) ([]*node.Node, error) {
	newNodes := make(map[string]*node.Node)
	for _, node := range nodes {
		if err := r.handleNode(node, newNodes); err != nil {
			return nil, errors.Wrap(err, "handle node "+node.Name)
		}
	}

	for _, newNode := range newNodes {
		nodes = append(nodes, newNode)
	}

	return nodes, nil
}

func (r *Remote) handleNode(n *node.Node, newNodes map[string]*node.Node) error {
	path := filepath.Join(r.root, n.Name)
	data, err := afero.ReadFile(r.fs, path)
	if err != nil {
		return errors.Wrap(err, "read file "+path)
	}
	logrus.Debugf("read file %s", path)

	includes, err := includes.Parse(data)
	if err != nil {
		return errors.Wrap(err, "parse includes")
	}

	depsAdded := make(map[string]bool)
	for _, include := range includes {
		depNode, err := r.resolveInclude(include)
		if err != nil {
			return errors.Wrap(err, "resolve include path "+include)
		} else if depNode == nil {
			continue
		}

		logrus.Debugf("resolved include %s to %s", include, depNode)
		if _, ok := depsAdded[depNode.Name]; !ok {
			logrus.Debugf("adding %s dependency to %s", depNode, n)
			n.Dependencies = append(n.Dependencies, depNode)
			depsAdded[depNode.Name] = true
		}

		newNodes[depNode.Name] = depNode
	}

	return nil
}

func (r *Remote) resolveInclude(include string) (*node.Node, error) {
	logrus.Debugf("resolving remote include %s", include)
	d := findDep(include)
	if d == nil {
		return nil, nil
	}

	destDir := filepath.Join(r.cache, "dependencies", d.name)

	if exists, err := afero.Exists(r.fs, destDir); err != nil {
		return nil, errors.Wrap(err, "exists")
	} else if exists {
		logrus.Debugf("already downloaded")
	} else {
		if err := r.fs.MkdirAll(destDir, 0755); err != nil {
			return nil, errors.Wrap(err, "mkdir all")
		}
		logrus.Debugf("downloading to dir %s", destDir)

		if err := r.d.Download(r.fs, destDir, d.url, d.sha256); err != nil {
			return nil, errors.Wrap(err, "download")
		}
	}

	n := &node.Node{
		Name: d.name,
		//Sources:      prependDir(d.sources, destDir),
		//Headers:      prependDir(d.headers, destDir),
		//IncludePaths: prependDir(d.includePaths, destDir),
	}

	return n, nil
}

func prependDir(ss []string, dir string) []string {
	ssNew := make([]string, len(ss))
	for i, s := range ss {
		ssNew[i] = filepath.Join(dir, s)
	}
	return ssNew
}
