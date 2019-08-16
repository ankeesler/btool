// Package btool provides fundamental entities that can be used to perform
// btool domain work.
package btool

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/ankeesler/btool/node/pipeline"
	"github.com/ankeesler/btool/node/pipeline/handlers"
	"github.com/ankeesler/btool/node/pipeline/handlers/resolverfactory"
	"github.com/ankeesler/btool/node/pipeline/handlers/store"
	registrypkg "github.com/ankeesler/btool/node/registry"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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

	rootAbs, err := filepath.Abs(cfg.Root)
	if err != nil {
		return errors.Wrap(err, "abs")
	}
	project := filepath.Base(rootAbs)
	projectDir := s.ProjectDir(project)

	logrus.Debugf("root: %s, project: %s", rootAbs, project)

	if err := os.MkdirAll(filepath.Dir(projectDir), 0755); err != nil {
		return errors.Wrap(err, "mkdir all")
	}

	if err := os.Symlink(rootAbs, projectDir); err != nil {
		return errors.Wrap(err, "symlink")
	}

	target := filepath.Join(projectDir, cfg.Target)

	ctx := pipeline.NewCtx()
	p := pipeline.New(ctx)

	rf := resolverfactory.New(
		cfg.CompilerC,
		cfg.CompilerCC,
		cfg.Archiver,
		cfg.Linker,
	)

	rhs, err := createRegistryHandlers(fs, s, rf, cfg.Registries)
	if err != nil {
		return errors.Wrap(err, "create registry handlers")
	}
	p.Handlers(rhs...)

	p.Handlers(
		handlers.NewFS(fs, projectDir),
		handlers.NewPrint(os.Stdout),
		handlers.NewObject(s, rf, project, target),
		handlers.NewExecutable(s, rf, project, target),
		handlers.NewSortAlpha(),
		handlers.NewResolve(fs, target),
	)

	if err := p.Run(); err != nil {
		return errors.Wrap(err, "pipeline run")
	}

	return nil
}

func createRegistryHandlers(
	fs afero.Fs,
	s handlers.Store,
	rf handlers.ResolverFactory,
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

		hs = append(hs, handlers.NewRegistry(fs, s, rf, r))
	}

	return hs, nil
}
