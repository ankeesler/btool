// Package btool provides fundamental entities that can be used to perform
// btool domain work.
package btool

import (
	"github.com/ankeesler/btool/node/pipeline"
	"github.com/ankeesler/btool/node/pipeline/handlers"
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
	).Linker(
		cfg.Linker,
	).Build()

	p := pipeline.New(
		ctx,
		handlers.NewFS(fs),
		handlers.NewDepsLocal(fs),
		handlers.NewObject(),
		handlers.NewExecutable(),
		handlers.NewSortAlpha(),
		handlers.NewResolve(fs),
	)
	if err := p.Run(); err != nil {
		return errors.Wrap(err, "pipeline run")
	}

	return nil
}
