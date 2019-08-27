// Package gaggler provides a type that can retrieve a registry.Gaggle.
package gaggler

import "github.com/ankeesler/btool/registry"

// Gaggler is a type that holds a registry.Gaggle.
type Gaggler struct {
	g *registry.Gaggle
}

func newGaggler(g *registry.Gaggle) *Gaggler {
	return &Gaggler{
		g: g,
	}
}

// Gaggle returns the registry.Gaggle that this Gaggler is holding.
func (g *Gaggler) Gaggle() (*registry.Gaggle, error) {
	return g.g, nil
}
