// Package scanner provides a prodcon.Producer that can produce node.Node's by
// walking a filesystem.
package scanner

import (
	collector "github.com/ankeesler/btool/collector0"
	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Walker

// Walker is a type that can walk a filesystem and return a list of paths.
type Walker interface {
	Walk(string, []string) ([]string, error)
}

// Scanner is a type that can Produce() node.Node's from a filesystem.
type Scanner struct {
	w    Walker
	root string
	exts []string
}

// New creates a new Scanner to run from a root on files with the provided
// extensions.
func New(w Walker, root string, exts []string) *Scanner {
	return &Scanner{
		w:    w,
		root: root,
		exts: exts,
	}
}

func (s *Scanner) Produce(store *collector.Store) error {
	files, err := s.w.Walk(s.root, s.exts)
	if err != nil {
		return errors.Wrap(err, "walk")
	}

	for _, file := range files {
		store.Add(node.New(file))
	}

	return nil
}
