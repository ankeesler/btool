// Package btool provides fundamental entities that can be used to perform
// btool domain work.
package btool

import (
	"path/filepath"

	"github.com/ankeesler/btool/builder"
	"github.com/ankeesler/btool/builder/currenter"
	"github.com/ankeesler/btool/collector"
	"github.com/ankeesler/btool/collector/resolverfactory"
	"github.com/ankeesler/btool/collector/scanner"
	"github.com/ankeesler/btool/collector/scanner/includeser"
	"github.com/ankeesler/btool/collector/scanner/nodestore"
	"github.com/ankeesler/btool/collector/sorter"
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

	fs := afero.NewOsFs()
	ns := nodestore.New(ui)
	i := includeser.New()
	rf := resolverfactory.New(
		cfg.CompilerC,
		cfg.CompilerCC,
		cfg.Archiver,
		cfg.LinkerC,
		cfg.LinkerCC,
	)
	scanner := scanner.New(fs, cfg.Root, ns, i, rf)
	sorter := sorter.New()
	target := filepath.Join(cfg.Root, cfg.Target)
	c := collector.New(scanner, sorter, target)

	targetN, err := c.Collect()
	if err != nil {
		return errors.Wrap(err, "collect")
	}

	cur := currenter.New()
	b := builder.New(cur, ui)
	if err := b.Build(targetN); err != nil {
		return errors.Wrap(err, "build")
	}

	return nil
}
