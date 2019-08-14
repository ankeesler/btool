// Package btool provides fundamental entities that can be used to perform
// btool domain work.
package btool

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/ankeesler/btool/node/pipeline"
	"github.com/ankeesler/btool/node/pipeline/handlers"
	registrypkg "github.com/ankeesler/btool/node/registry"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// Cfg is a configuration struct provided to a Build call.
//
// Callers should set all fields.
type Cfg struct {
	Root   string
	Cache  string
	Target string

	CompilerC  string
	CompilerCC string
	Archiver   string
	Linker     string

	Registries []string
}

// Run will run a btool build and produce a target.
func Run(cfg *Cfg) error {
	fs := afero.NewOsFs()

	ctx := pipeline.NewCtxBuilder().Root(
		cfg.Root,
	).Cache(
		cfg.Cache,
	).Target(
		cfg.Target,
	).CompilerC(
		cfg.CompilerC,
	).CompilerCC(
		cfg.CompilerCC,
	).Archiver(
		cfg.Archiver,
	).Linker(
		cfg.Linker,
	).Build()

	p := pipeline.New(ctx)

	if err := addRegistryHandlers(cfg.Cache, p, fs, cfg.Registries); err != nil {
		return errors.Wrap(err, "add registry handlers")
	}

	p.Handler(
		handlers.NewFS(fs),
	).Handler(
		handlers.NewObject(),
	).Handler(
		handlers.NewExecutable(),
	).Handler(
		handlers.NewSortAlpha(),
	).Handler(
		handlers.NewResolve(fs),
	)

	if err := p.Run(); err != nil {
		return errors.Wrap(err, "pipeline run")
	}

	return nil
}

func addRegistryHandlers(
	cache string,
	p *pipeline.Pipeline,
	fs afero.Fs,
	registries []string,
) error {
	for _, registry := range registries {
		url, err := url.Parse(registry)
		if err != nil {
			return errors.Wrap(err, "url parse")
		}

		var r handlers.Registry
		switch url.Scheme {
		case "http", "https":
			c := &http.Client{}
			r = registrypkg.NewHTTPRegistry(registry, c)
		case "file":
			r, err = registrypkg.CreateFSRegistry(fs, url.Path)
		default:
			r, err = registrypkg.CreateFSRegistry(fs, registry)
		}

		registryStorePath := makeRegistryStorePath(cache, registry)
		projectsStorePath := makeProjectsStorePath(cache)
		h := handlers.NewRegistry(fs, r, registryStorePath, projectsStorePath)
		p.Handler(h)
	}
	return nil
}

func makeRegistryStorePath(cache, registry string) string {
	r := base64.StdEncoding.EncodeToString([]byte(registry))
	return filepath.Join(cache, "registries", r)
}

func makeProjectsStorePath(cache string) string {
	return filepath.Join(cache, "projects")
}
