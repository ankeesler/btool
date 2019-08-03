package builder2

import (
	"errors"

	"github.com/ankeesler/btool/scanner/graph"
)

func Build(g *graph.Graph, target string) error {
	b := g.Find(target)
	if b == nil {
		return errors.New("unknown target: " + target)
	}

	dependencies := g.Dependencies(target)
	for dependency := range dependencies {
		if err := Build(g, dependency.path); err != nil {
			return err
		}
	}

	if err := b.Resolve(g, path); err != nil {
		return errors.Wrap(err, "resolve")
	}

	return nil
}
