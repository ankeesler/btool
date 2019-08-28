// Package gaggler provides a type that can retrieve a registry.Gaggle.
package gaggler

import "github.com/ankeesler/btool/registry"

// Gaggler is a type that holds a registry.Gaggle.
type Gaggler struct {
	g    *registry.Gaggle
	root string
}

func newGaggler(g *registry.Gaggle, root string) *Gaggler {
	return &Gaggler{
		g:    g,
		root: root,
	}
}

// Gaggle returns the registry.Gaggle that this Gaggler is holding.
func (g *Gaggler) Gaggle() *registry.Gaggle {
	return g.g
}

func (g *Gaggler) Root() string {
	return g.root
}
