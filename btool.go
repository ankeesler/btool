// Package btool provides fundamental entities that can be used to perform
// btool domain work.
package btool

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

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
	Output string

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

	info, err := os.Stat(projectDir)
	logrus.Debugf("examining projectDir %s (%s)", projectDir, info)
	if err != nil {
		if !os.IsNotExist(err) {
			return errors.Wrap(err, "stat")
		}
	} else if (info.Mode() & os.ModeSymlink) != 0 {
		return fmt.Errorf("expected %s to be symlink (%s)", projectDir, info)
	} else if err := os.Remove(projectDir); err != nil {
		return errors.Wrap(err, "remote")
	}
	if err := os.Symlink(rootAbs, projectDir); err != nil {
		return errors.Wrap(err, "symlink")
	}
	logrus.Debugf("symlinked %s to %s", projectDir, rootAbs)

	var target string
	if strings.HasPrefix(cfg.Target, cfg.Cache) {
		target = cfg.Target
	} else {
		target = filepath.Join(projectDir, cfg.Target)
	}

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
		handlers.NewObject(s, rf, project, target),
		handlers.NewExecutable(s, rf, project, target),
		handlers.NewSymlink(rf, cfg.Output, target),
		handlers.NewSortAlpha(),
		handlers.NewResolve(fs, cfg.Output),
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
			logrus.Debugf("creating http registry from %s", registry)
		case "file":
			r, err = registrypkg.CreateFSRegistry(fs, url.Path, registry)
			logrus.Debugf("creating fs registry from %s", url.Path)
		default:
			r, err = registrypkg.CreateFSRegistry(fs, registry, registry)
			logrus.Debugf("creating fs registry from %s", registry)
		}

		hs = append(hs, handlers.NewRegistry(fs, s, rf, r))
	}

	return hs, nil
}
