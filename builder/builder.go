// Package builder is in charge of walking the project graph to compile and link
// stuff.
package builder

import (
	"path/filepath"

	"github.com/ankeesler/btool/config"
	"github.com/spf13/afero"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Toolchain

type Toolchain interface {
	CompileC(output, input string, includeDirs []string) error
	CompileCC(output, input string, includeDirs []string) error
	Link(output string, inputs []string) error
}

type Builder struct {
	fs     afero.Fs
	config *config.Config
	t      Toolchain
}

func New(fs afero.Fs, config *config.Config, t Toolchain) *Builder {
	return &Builder{
		fs:     fs,
		config: config,
		t:      t,
	}
}

func (b *Builder) objectsDir() string {
	return filepath.Join(
		b.config.Cache,
		"objects",
	)
}

func (b *Builder) binariesDir() string {
	return filepath.Join(
		b.config.Cache,
		"binaries",
	)
}
