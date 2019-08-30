// Package btool provides a simple Run() function that clients can call to
// run the btool application. This Run() function performs the necessary
// dependency injection so that callers do not have to.
//   ...
//   err := Run(&Cfg{
//     ...
//   })
//   ...

package btool

import (
	"path/filepath"

	"github.com/ankeesler/btool/app"
	"github.com/ankeesler/btool/builder"
	"github.com/ankeesler/btool/builder/currenter"
	"github.com/ankeesler/btool/cleaner"
	"github.com/ankeesler/btool/collector"
	"github.com/ankeesler/btool/collector/registry"
	"github.com/ankeesler/btool/collector/registry/gaggle"
	"github.com/ankeesler/btool/collector/resolverfactory"
	"github.com/ankeesler/btool/collector/scanner"
	"github.com/ankeesler/btool/collector/scanner/includeser"
	"github.com/ankeesler/btool/collector/sorter"
	"github.com/ankeesler/btool/node"
	registrypkg "github.com/ankeesler/btool/registry"
	"github.com/ankeesler/btool/ui"
	"github.com/spf13/afero"
)

// Cfg is a configuration struct provided to a Build call.
//
// Callers should set all fields.
type Cfg struct {
	Root   string
	Cache  string
	Target string

	DryRun bool
	Clean  bool

	CompilerC  string
	CompilerCC string
	Archiver   string
	LinkerC    string
	LinkerCC   string

	Registry string

	Quiet bool
}

type collectorCreator struct {
	ctx    *collector.Ctx
	cinics []collector.CollectiniCreator
}

func (cc *collectorCreator) Create() (app.Collector, error) {
	return collector.NewCreator(cc.ctx, cc.cinics).Create()
}

type registryClientCreator struct {
	url string
}

func (rc *registryClientCreator) Create() (registry.Client, error) {
	return registrypkg.NewCreator(rc.url).Create()
}

type registryCollectiniCreator struct {
	fs    afero.Fs
	rcc   *registryClientCreator
	cache string
	gc    *gaggle.Collector
}

func (cc *registryCollectiniCreator) Create() (collector.Collectini, error) {
	return registry.NewCreator(
		cc.fs,
		cc.rcc,
		cc.cache,
		cc.gc,
	).Create()
}

// Run will run a btool invocation and produce a target.
//
// Under the hood, Run creates the dependencies for a Btool struct via the
// provided Cfg, passes those dependencies to New(), and calls Run() on the
// returned Btool struct.
func Run(cfg *Cfg) error {
	ui := ui.New(cfg.Quiet)

	fs := afero.NewOsFs()
	ns := collector.NewNodeStore(ui)
	i := includeser.New(fs)
	rf := resolverfactory.New(
		cfg.CompilerC,
		cfg.CompilerCC,
		cfg.Archiver,
		cfg.LinkerC,
		cfg.LinkerCC,
	)

	ctx := collector.NewCtx(ns, rf)

	cinics := []collector.CollectiniCreator{
		&registryCollectiniCreator{
			fs:    fs,
			rcc:   &registryClientCreator{url: cfg.Registry},
			cache: cfg.Cache,
			gc:    gaggle.New(),
		},
		collector.NewCollectiniAccessor(
			scanner.New(fs, cfg.Root, i),
		),
		collector.NewCollectiniAccessor(
			sorter.New(),
		),
	}
	cc := &collectorCreator{ctx: ctx, cinics: cinics}
	cleaner := cleaner.New(fs, ui)
	builder := builder.New(cfg.DryRun, currenter.New(), ui)
	a := app.New(cc, cleaner, builder)

	target := filepath.Join(cfg.Root, cfg.Target)
	targetN := node.New(target)
	return a.Run(targetN, cfg.Clean, cfg.DryRun)
}
