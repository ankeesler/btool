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
	"github.com/ankeesler/btool/node/pipeline/handlers/store"
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
	s := store.New(cfg.Cache)

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

	rhs, err := createRegistryHandlers(fs, s, cfg.Registries)
	if err != nil {
		return errors.Wrap(err, "create registry handlers")
	}
	for _, rh := range rhs {
		p.Handler(rh)
	}

	p.Handler(
		handlers.NewFS(fs),
	).Handler(
		handlers.NewObject(s),
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

func createRegistryHandlers(
	fs afero.Fs,
	s handlers.Store,
	registries []string,
) ([]pipeline.Handler, error) {
	hs := make([]pipeline.Handler, 0)

	for _, registry := range registries {
		url, err := url.Parse(registry)
		if err != nil {
			return nil, errors.Wrap(err, "url parse")
		}

		var r handlers.Registry
		switch url.Scheme {
		case "http", "https":
			c := &http.Client{}
			r = registrypkg.NewHTTPRegistry(registry, c)
		case "file":
			r, err = registrypkg.CreateFSRegistry(fs, url.Path, registry)
		default:
			r, err = registrypkg.CreateFSRegistry(fs, registry, registry)
		}

		hs = append(hs, handlers.NewRegistry(fs, s, r))
	}

	return hs, nil
}

func makeRegistryStorePath(cache, registry string) string {
	r := base64.StdEncoding.EncodeToString([]byte(registry))
	return filepath.Join(cache, "registries", r)
}

func makeProjectsStorePath(cache string) string {
	return filepath.Join(cache, "projects")
}
