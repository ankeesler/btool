// Package linker provides a node.Handler that will link a target.
package linker

import (
	"github.com/ankeesler/btool/node"
	"github.com/spf13/afero"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . L

// L is an actual linker client.
type L interface {
	Link(output string, inputs []string) error
}

type Linker struct {
	l      L
	fs     afero.Fs
	cache  string
	target string
}

func New(l L, fs afero.Fs, cache, target string) *Linker {
	return &Linker{
		l:      l,
		fs:     fs,
		cache:  cache,
		target: target,
	}
}

func (l *Linker) Handler(nodes []*node.Node) ([]*node.Node, error) {
	// Find target.
	// Collect objects.
	// Invoke linker.
	return nodes, nil
}
