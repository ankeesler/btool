// Package builder is in charge of walking the project graph to compile and link
// stuff.
package builder

import (
	"path/filepath"

	"github.com/ankeesler/btool/config"
	"github.com/spf13/afero"
)

type Compiler interface {
	CompileC(output, input, root string) error
	CompileCC(output, input, root string) error
}

type Linker interface {
	Link(output string, inputs []string) error
}

type Builder struct {
	fs     afero.Fs
	config *config.Config
	c      Compiler
	l      Linker
}

func New(fs afero.Fs, config *config.Config, c Compiler, l Linker) *Builder {
	return &Builder{
		fs:     fs,
		config: config,
		c:      c,
		l:      l,
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
