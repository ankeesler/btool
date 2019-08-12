package handlers

import (
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/pipeline"
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
	// registrypkg.IndexFile.Path. If no such object exists, nil should be returned.
	// If there is an error, an error should be returned.
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
		project := getProject(file.Path)
		path := cacheDownloadPath(ctx, project, file.SHA256)
		logrus.Debugf("considering %s", path)

		var nodes []*registrypkg.Node
		if exists, err := afero.Exists(r.fs, path); err != nil {
			return errors.Wrap(err, "exists")
		} else if !exists {
			logrus.Debugf("does not exist")

			nodes, err = r.r.Nodes(file.Path)
			if err != nil {
				return errors.Wrap(err, "nodes")
			} else if nodes == nil {
				return errors.New("unknown nodes at path: " + file.Path)
			}
		} else {
			logrus.Debugf("in cache")

			data, err := afero.ReadFile(r.fs, path)
			if err != nil {
				return errors.Wrap(err, "read file")
			}

			nodes = make([]*registrypkg.Node, 0)
			if err := yaml.Unmarshal(data, &nodes); err != nil {
				return errors.Wrap(err, "unmarshal")
			}
		}

		for _, n := range nodes {
			nN, err := r.d.Decode(n)
			if err != nil {
				return errors.Wrap(err, "decode")
			}

			logrus.Debugf("decoded %s to %s", n, n)
			ctx.Nodes = append(ctx.Nodes, nN)
		}
	}

	return nil
}

func (r *registry) Name() string { return "registry" }
