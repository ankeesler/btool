// Package builder is in charge of walking the project graph to compile and link
// stuff.
package builder

import (
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
	fs    afero.Fs
	root  string
	store string
	c     Compiler
	l     Linker
}

func New(fs afero.Fs, root, store string, c Compiler, l Linker) *Builder {
	return &Builder{
		fs:    fs,
		root:  root,
		store: store,
		c:     c,
		l:     l,
	}
}
