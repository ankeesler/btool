// Package btool provides fundamental entities that can be used to perform
// btool domain work.
package btool

import (
	"github.com/ankeesler/btool/ui"
	"github.com/pkg/errors"
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
	LinkerC    string
	LinkerCC   string

	Registries []string

	Quiet bool
}

// Run will run a btool build and produce a target.
// Run calls Collect() and then Build().
func Run(cfg *Cfg) error {
	ui := ui.New(cfg.Quiet)

	targetN, err := Collect(cfg, ui)
	if err != nil {
		return errors.Wrap(err, "collect")
	}

	if err := Build(targetN, ui); err != nil {
		return errors.Wrap(err, "build")
	}

	return nil
}
