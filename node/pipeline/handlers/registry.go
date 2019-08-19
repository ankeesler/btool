package handlers

import (
	"fmt"
	"path/filepath"

	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/pipeline"
	registrypkg "github.com/ankeesler/btool/node/registry"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
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
		log.Debugf("considering %s", gaggleFile)
		if exists, err := afero.Exists(r.fs, gaggleFile); err != nil {
			return errors.Wrap(err, "exists")
		} else if !exists {
			log.Debugf("does not exist")

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
			Project     string   `mapstructure:"project"`
			IncludeDirs []string `mapstructure:"includeDirs"`
		}{}
		if err := mapstructure.Decode(gaggle.Metadata, &metadata); err != nil {
			return errors.Wrap(err, "decode")
		}
		log.Debugf("metadata: %+v", metadata)

		projectDir := r.s.ProjectDir(metadata.Project)
		includeDirs := prependDir(metadata.IncludeDirs, projectDir)

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

			nodeR, err := r.newResolver(n.Resolver, includeDirs, projectDir)
			if err != nil {
				return errors.Wrap(err, "new resolver")
			}
			nN.Resolver = nodeR

			log.Debugf("decoded %s to %s", n, nN)
			ctx.Nodes = append(ctx.Nodes, nN)
		}
	}

	return nil
}

func (r *registry) String() string { return "registry" }

func (r *registry) newResolver(
	registryR registrypkg.Resolver,
	includeDirs []string,
	projectDir string,
) (node.Resolver, error) {
	name := registryR.Name
	config := registryR.Config

	var nodeR node.Resolver
	var err error
	switch name {
	case "compileC":
		nodeR = r.rf.NewCompileC(includeDirs)
	case "compileCC":
		nodeR = r.rf.NewCompileCC(includeDirs)
	case "archive":
		nodeR = r.rf.NewArchive()
	case "linkC":
		nodeR = r.rf.NewLinkC()
	case "linkCC":
		nodeR = r.rf.NewLinkCC()
	case "symlink":
		nodeR = r.rf.NewSymlink()
	case "unzip":
		nodeR = r.rf.NewUnzip(projectDir)
	case "download":
		nodeR, err = r.createDownload(config)
		if err != nil {
			err = errors.Wrap(err, "create download")
		}
	default:
		err = fmt.Errorf("unknown resolver: %s", name)
	}

	return nodeR, err
}

func (r *registry) createDownload(
	config map[string]interface{},
) (node.Resolver, error) {
	cfg := struct {
		URL    string
		SHA256 string
	}{}
	if err := mapstructure.Decode(config, &cfg); err != nil {
		return nil, errors.Wrap(err, "decode")
	}

	return r.rf.NewDownload(cfg.URL, cfg.SHA256), nil
}

func prependDir(dirs []string, dir string) []string {
	for i := range dirs {
		dirs[i] = filepath.Join(dir, dirs[i])
	}
	return dirs
}
