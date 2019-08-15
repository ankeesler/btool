package handlers

import (
	"fmt"
	"path/filepath"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/pipeline"
	registrypkg "github.com/ankeesler/btool/node/registry"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Registry

// Registry is an object that can retrieve registry.Gaggle's.
type Registry interface {
	// Index should return the registrypkg.Index associated with this particular
	// Registry. If any error occurs, an error should be returned.
	Index() (*registrypkg.Index, error)
	// Gaggle should return the registrypkg.Gaggle associated with the provided
	// registrypkg.IndexFile.Path. If any error occurs, an error should be returned.
	// If no registrypkg.Gaggle exists for the provided string, then nil, nil should
	// be returned.
	Gaggle(string) (*registrypkg.Gaggle, error)
}

type registry struct {
	fs afero.Fs
	s  Store
	rf ResolverFactory
	r  Registry
}

// NewRegistry returns a pipeline.Handler that retrieves node.Node's from a
// Registry.
func NewRegistry(
	fs afero.Fs,
	s Store,
	rf ResolverFactory,
	r Registry,
) pipeline.Handler {
	return &registry{
		fs: fs,
		s:  s,
		rf: rf,
		r:  r,
	}
}

func (r *registry) Handle(ctx *pipeline.Ctx) error {
	i, err := r.r.Index()
	if err != nil {
		return errors.Wrap(err, "index")
	}

	registryDir := r.s.RegistryDir(i.Name)

	for _, file := range i.Files {

		gaggleFile := filepath.Join(registryDir, file.SHA256+".yml")
		gaggle := new(registrypkg.Gaggle)
		logrus.Debugf("considering %s", gaggleFile)
		if exists, err := afero.Exists(r.fs, gaggleFile); err != nil {
			return errors.Wrap(err, "exists")
		} else if !exists {
			logrus.Debugf("does not exist")

			gaggle, err = r.r.Gaggle(file.Path)
			if err != nil {
				return errors.Wrap(err, "gaggle")
			} else if gaggle == nil {
				return errors.New("unknown gaggle at path: " + file.Path)
			}

			gaggleData, err := yaml.Marshal(&gaggle)
			if err != nil {
				return errors.Wrap(err, "marshal")
			}

			if err := r.fs.MkdirAll(filepath.Dir(gaggleFile), 0755); err != nil {
				return errors.Wrap(err, "mkdir all")
			}

			if err := afero.WriteFile(r.fs, gaggleFile, gaggleData, 0644); err != nil {
				return errors.Wrap(err, "write file")
			}
		} else {
			data, err := afero.ReadFile(r.fs, gaggleFile)
			if err != nil {
				return errors.Wrap(err, "read file")
			}

			if err := yaml.Unmarshal(data, &gaggle); err != nil {
				return errors.Wrap(err, "unmarshal")
			}
		}

		metadata := struct {
			Project string
		}{}
		if err := mapstructure.Decode(gaggle.Metadata, &metadata); err != nil {
			return errors.Wrap(err, "decode")
		}

		projectDir := r.s.ProjectDir(metadata.Project)

		for _, n := range gaggle.Nodes {
			nN := node.New(filepath.Join(projectDir, n.Name))
			for _, d := range n.Dependencies {
				dName := filepath.Join(projectDir, d)

				var dN *node.Node
				if d == "$this" {
					// TODO: test me.
					dN = node.New(gaggleFile)
				} else {
					dN = node.Find(dName, ctx.Nodes)
				}

				if dN == nil {
					return fmt.Errorf("cannot find dependency %s/%s of %s", d, dName, n)
				}
				nN.Dependency(dN)
			}

			n.Resolver.Config["outputDir"] = projectDir
			r, err := r.rf.NewResolver(n.Resolver.Name, n.Resolver.Config)
			if err != nil {
				return errors.Wrap(err, "new resolver")
			}
			nN.Resolver = r

			logrus.Debugf("decoded %s to %s", n, nN)
			ctx.Nodes = append(ctx.Nodes, nN)
		}
	}

	return nil
}

func (r *registry) Name() string { return "registry" }
