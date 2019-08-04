package btool

import (
	"github.com/ankeesler/btool/pipeline"
	"github.com/ankeesler/btool/pipeline/handlers"
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
	Linker     string
}

// Build will run a btool build and produce a target.
func Build(cfg *Cfg) error {
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
	).Linker(
		cfg.Linker,
	).Build()

	p := pipeline.New(
		ctx,
		handlers.NewFS(fs),
		handlers.NewDepsLocal(fs),
		handlers.NewObject(),
		handlers.NewExecutable(),
		handlers.NewResolve(),
	)
	if err := p.Run(); err != nil {
		return errors.Wrap(err, "pipeline run")
	}

	return nil
}
