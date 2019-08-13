package handlers

import (
	"fmt"
	"path/filepath"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/pipeline"
	"github.com/ankeesler/btool/node/pipeline/resolvermapper"
	registrypkg "github.com/ankeesler/btool/node/registry"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Registry

// Registry is an object that can retrieve registry.Node's.
type Registry interface {
	// Index should return the registrypkg.Index associated with this particular
	// Registry. If any error occurs, an error should be returned.
	Index() (*registrypkg.Index, error)
	// Nodes should return the registrypkg.Node's associated with the provided
	// registrypkg.IndexFile.Path. If any error occurs, an error should be returned.
	// If no registrypkg.Node's exist for the provided string, then an empty
	// slice should be returned.
	Nodes(string) ([]*registrypkg.Node, error)
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Decoder

// Decoder is an object that can change a registry.Node into a node.Node.
type Decoder interface {
	Decode(*registrypkg.Node) (*node.Node, error)
}

type registry struct {
	fs afero.Fs
	r  Registry
	d  Decoder
}

// NewRegistry returns a pipeline.Handler that retrieves node.Node's from a
// Registry.
func NewRegistry(fs afero.Fs, r Registry, d Decoder) pipeline.Handler {
	return &registry{
		fs: fs,
		r:  r,
		d:  d,
	}
}

func (r *registry) Handle(ctx *pipeline.Ctx) error {
	i, err := r.r.Index()
	if err != nil {
		return errors.Wrap(err, "index")
	}

	for _, file := range i.Files {
		cachePath := filepath.Join(ctx.KV[pipeline.CtxCache], file.SHA256)

		nodesFile := filepath.Join(cachePath, "nodes.yml")
		nodes := make([]*registrypkg.Node, 0)
		logrus.Debugf("considering %s", nodesFile)
		if exists, err := afero.Exists(r.fs, nodesFile); err != nil {
			return errors.Wrap(err, "exists")
		} else if !exists {
			logrus.Debugf("does not exist")

			nodes, err = r.r.Nodes(file.Path)
			if err != nil {
				return errors.Wrap(err, "nodes")
			} else if nodes == nil {
				return errors.New("unknown nodes at path: " + file.Path)
			}

			nodesData, err := yaml.Marshal(nodes)
			if err != nil {
				return errors.Wrap(err, "marshal")
			}

			if err := r.fs.MkdirAll(filepath.Dir(nodesFile), 0755); err != nil {
				return errors.Wrap(err, "mkdir all")
			}

			if err := afero.WriteFile(r.fs, nodesFile, nodesData, 0644); err != nil {
				return errors.Wrap(err, "write file")
			}
		} else {
			data, err := afero.ReadFile(r.fs, nodesFile)
			if err != nil {
				return errors.Wrap(err, "read file")
			}

			if err := yaml.Unmarshal(data, &nodes); err != nil {
				return errors.Wrap(err, "unmarshal")
			}
		}

		nodesDir := filepath.Join(cachePath, "nodes")
		for _, n := range nodes {
			nN := node.New(filepath.Join(nodesDir, n.Name))
			for _, d := range n.Dependencies {
				dName := filepath.Join(nodesDir, d)
				dN := node.Find(dName, ctx.Nodes)
				if dN == nil {
					return fmt.Errorf("cannot find dependency %s/%s of %s", d, dName, n)
				}
			}
			rm := resolvermapper.New(ctx)
			nN.Resolver, err = rm.Map(n.Resolver.Name, n.Resolver.Config)
			if err != nil {
				return errors.Wrap(err, "map")
			}

			logrus.Debugf("decoded %s to %s", n, nN)
			ctx.Nodes = append(ctx.Nodes, nN)
		}
	}

	return nil
}

func (r *registry) Name() string { return "registry" }
