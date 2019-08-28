package btool

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
	//ui := ui.New(cfg.Quiet)
	//
	//fs := afero.NewOsFs()
	//ns := collector.NewNodeStore(ui)
	//i := includeser.New(fs)
	//rf := resolverfactory.New(
	//	cfg.CompilerC,
	//	cfg.CompilerCC,
	//	cfg.Archiver,
	//	cfg.LinkerC,
	//	cfg.LinkerCC,
	//)
	//scanner := scanner.New(fs, cfg.Root, i)
	//sorter := sorter.New()
	//
	//ctx := collector.NewCtx(ns, rf)
	//collector := collector.New(ctx, scanner, sorter)
	//cleaner := cleaner.New(fs, ui)
	//builder := builder.New(cfg.DryRun, currenter.New(), ui)
	//
	//target := filepath.Join(cfg.Root, cfg.Target)
	//targetN := node.New(target)
	//b := New(collector, cleaner, builder)
	//return b.Run(targetN, cfg.Clean, cfg.DryRun)
	return nil
}
