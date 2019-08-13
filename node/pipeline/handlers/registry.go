package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/pipeline"
	registrypkg "github.com/ankeesler/btool/node/registry"
	"github.com/ankeesler/btool/node/resolvers"
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
	r  Registry
}

// NewRegistry returns a pipeline.Handler that retrieves node.Node's from a
// Registry.
func NewRegistry(fs afero.Fs, r Registry) pipeline.Handler {
	return &registry{
		fs: fs,
		r:  r,
	}
}

func (r *registry) Handle(ctx *pipeline.Ctx) error {
	i, err := r.r.Index()
	if err != nil {
		return errors.Wrap(err, "index")
	}

	for _, file := range i.Files {
		cachePath := filepath.Join(ctx.KV[pipeline.CtxCache], file.SHA256)

		gaggleFile := filepath.Join(cachePath, "gaggle.yml")
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

		gaggleDir := filepath.Join(cachePath, "gaggle")
		for _, n := range gaggle.Nodes {
			nN := node.New(filepath.Join(gaggleDir, n.Name))
			for _, d := range n.Dependencies {
				dName := filepath.Join(gaggleDir, d)
				dN := node.Find(dName, ctx.Nodes)
				if dN == nil {
					return fmt.Errorf("cannot find dependency %s/%s of %s", d, dName, n)
				}
			}

			if err := setResolver(ctx, nN, n.Resolver.Name, n.Resolver.Config); err != nil {
				return errors.Wrap(err, "map")
			}

			logrus.Debugf("decoded %s to %s", n, nN)
			ctx.Nodes = append(ctx.Nodes, nN)
		}
	}

	return nil
}

func (r *registry) Name() string { return "registry" }

func setResolver(
	ctx *pipeline.Ctx,
	n *node.Node,
	name string,
	config map[string]interface{},
) error {
	root := ctx.KV[pipeline.CtxRoot]
	cache := ctx.KV[pipeline.CtxCache]
	compilerC := ctx.KV[pipeline.CtxCompilerC]
	compilerCC := ctx.KV[pipeline.CtxCompilerCC]
	archiver := ctx.KV[pipeline.CtxArchiver]
	linker := ctx.KV[pipeline.CtxLinker]

	var r node.Resolver
	var err error
	switch name {
	case "compileC":
		r = resolvers.NewCompile(root, compilerC, []string{root})
	case "compileCC":
		r = resolvers.NewCompile(root, compilerCC, []string{root})
	case "archive":
		r = resolvers.NewArchive(root, archiver)
	case "link":
		r = resolvers.NewLink(root, linker)
	case "unzip":
		var err error
		r, err = createUnzip(config, cache)
		if err != nil {
			err = errors.Wrap(err, "create download")
		}
	case "download":
		var err error
		r, err = createDownload(config, cache)
		if err != nil {
			err = errors.Wrap(err, "create download")
		}
	default:
		err = fmt.Errorf("unknown resolver: %s", name)
	}

	if err != nil {
		return err
	}

	n.Resolver = r

	return nil
}

func createUnzip(
	config map[string]interface{},
	cache string,
) (node.Resolver, error) {
	outputDir := filepath.Join(cache, "unzip")
	return resolvers.NewUnzip(outputDir), nil
}

func createDownload(
	config map[string]interface{},
	cache string,
) (node.Resolver, error) {
	c := &http.Client{}

	cfg := struct {
		URL    string
		SHA256 string
	}{}
	if err := mapstructure.Decode(config, &cfg); err != nil {
		return nil, errors.Wrap(err, "decode")
	}

	outputFile := filepath.Join(cache, "download", cfg.SHA256)

	return resolvers.NewDownload(c, cfg.URL, cfg.SHA256, outputFile), nil
}
