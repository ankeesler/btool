// Package gaggler provides a type that can retrieve a registry.Gaggle.
package gaggler

import "github.com/ankeesler/btool/registry"

type Gaggler struct {
	g *registry.Gaggle
}

func newGaggler(g *registry.Gaggle) *Gaggler {
	return &Gaggler{
		g: g,
	}
}

func (g *Gaggler) Gaggle() (*registry.Gaggle, error) {
	return g.g, nil
}
