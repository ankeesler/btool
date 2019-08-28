package btool

import (
	"path/filepath"

	"github.com/ankeesler/btool/builder"
	"github.com/ankeesler/btool/builder/currenter"
	"github.com/ankeesler/btool/cleaner"
	"github.com/ankeesler/btool/collector"
	registrypkg "github.com/ankeesler/btool/collector"
	"github.com/ankeesler/btool/collector/registry"
	"github.com/ankeesler/btool/collector/registry/gaggle"
	"github.com/ankeesler/btool/collector/resolverfactory"
	"github.com/ankeesler/btool/collector/scanner"
	"github.com/ankeesler/btool/collector/scanner/includeser"
	"github.com/ankeesler/btool/collector/sorter"
	"github.com/ankeesler/btool/node"
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

	Registries []string

	Quiet bool
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
		collector.NewCollectiniAccessor(
			registry.NewCreator(
				fs,
				registrypkg.NewCreator(),
				cfg.Cache,
				gaggle.New(),
			),
		),
		collector.NewCollectiniAccessor(
			scanner.New(fs, cfg.Root, i),
		),
		collector.NewCollectiniAccessor(
			sorter.New(),
		),
	}
	cc := collector.NewCreator(ctx, cinics)
	ccf := func() (Collector, error) {
		return cc.Create()
	}
	cleaner := cleaner.New(fs, ui)
	builder := builder.New(cfg.DryRun, currenter.New(), ui)
	b := New(ccf, cleaner, builder)

	target := filepath.Join(cfg.Root, cfg.Target)
	targetN := node.New(target)
	return b.Run(targetN, cfg.Clean, cfg.DryRun)
}
