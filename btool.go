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
	"os"
	"path/filepath"

	"github.com/ankeesler/btool/app"
	"github.com/ankeesler/btool/app/builder"
	"github.com/ankeesler/btool/app/builder/currenter"
	"github.com/ankeesler/btool/app/cleaner"
	"github.com/ankeesler/btool/app/collector"
	"github.com/ankeesler/btool/app/collector/cc"
	"github.com/ankeesler/btool/app/collector/cc/includeser"
	"github.com/ankeesler/btool/app/collector/registry"
	"github.com/ankeesler/btool/app/collector/registry/clientcreator"
	"github.com/ankeesler/btool/app/collector/registry/gaggle"
	"github.com/ankeesler/btool/app/collector/resolverfactory"
	"github.com/ankeesler/btool/app/collector/scanner"
	"github.com/ankeesler/btool/app/collector/scanner/walker"
	"github.com/ankeesler/btool/app/lister"
	"github.com/ankeesler/btool/app/runner"
	"github.com/ankeesler/btool/app/watcher"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/ui"
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

	Clean  bool
	List   bool
	Run    bool
	Watch  bool
	DryRun bool

	CompilerC  string
	CompilerCC string
	Archiver   string
	LinkerC    string
	LinkerCC   string

	Registry string

	Quiet bool
}

type collectorCreator struct {
	fs afero.Fs
	rf *resolverfactory.ResolverFactory
	ui *ui.UI

	root     string
	cache    string
	targetN  *node.Node
	registry string
}

func (ccreator *collectorCreator) Create() (app.Collector, error) {
	rc := registry.NewCreator(
		ccreator.fs,
		clientcreator.New(ccreator.fs, ccreator.registry),
		ccreator.cache,
		gaggle.New(ccreator.rf),
	)
	r, err := rc.Create()
	if err != nil {
		return nil, errors.Wrap(err, "create registry")
	}
	s := scanner.New(
		walker.New(),
		ccreator.root,
		[]string{
			".c",
			".cc",
			".h",
		},
	)
	t := collector.NewTrivialProducer(ccreator.targetN)
	producers := []collector.Producer{
		r,
		s,
		t,
	}

	i := cc.NewIncludes(includeser.New(ccreator.fs))
	o := cc.NewObj(ccreator.rf)
	e := cc.NewExe(ccreator.rf)
	consumers := []collector.Consumer{
		i,
		o,
		e,
		ccreator.ui,
	}

	return collector.New(producers, consumers), nil
}

// Run will run a btool invocation and produce a target.
//
// Under the hood, Run creates the dependencies for an app.App struct via the
// provided Cfg, passes those dependencies to app.New(), and calls Run() on the
// returned app.App struct.
func Run(cfg *Cfg) error {
	ui := ui.New(cfg.Quiet)

	fs := afero.NewOsFs()
	rf := resolverfactory.New(
		cfg.CompilerC,
		cfg.CompilerCC,
		cfg.Archiver,
		cfg.LinkerC,
		cfg.LinkerCC,
	)

	target := filepath.Join(cfg.Root, cfg.Target)
	targetN := node.New(target)

	cc := &collectorCreator{
		fs: fs,
		rf: rf,
		ui: ui,

		root:     cfg.Root,
		cache:    cfg.Cache,
		targetN:  targetN,
		registry: cfg.Registry,
	}
	cleaner := cleaner.New(fs, ui)
	lister := lister.New(os.Stdout)
	builder := builder.New(cfg.DryRun, currenter.New(), ui)
	runner := runner.New(ui)
	watcher := watcher.New()
	a := app.New(cc, cleaner, lister, builder, runner, watcher)

	return a.Run(targetN, cfg.Clean, cfg.List, cfg.Run, cfg.Watch)
}
