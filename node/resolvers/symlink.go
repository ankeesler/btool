package resolvers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
)

type symlink struct {
}

// NewSymlink returns a node.Resolver that symlinks a node.Node to its
// dependency.
func NewSymlink() node.Resolver {
	return &symlink{}
}

func (s *symlink) Resolve(n *node.Node) error {
	if len(n.Dependencies) != 1 {
		return fmt.Errorf("expected 1 dependency, got %d", len(n.Dependencies))
	}

	if err := os.MkdirAll(filepath.Dir(n.Name), 0755); err != nil {
		return errors.Wrap(err, "mkdir all")
	}

	if _, err := os.Stat(n.Name); err != nil {
		if !os.IsNotExist(err) {
			return errors.Wrap(err, "stat")
		}
	} else if err := os.Remove(n.Name); err != nil {
		return errors.Wrap(err, "remove")
	}

	if err := os.Symlink(n.Dependencies[0].Name, n.Name); err != nil {
		return errors.Wrap(err, "symlink")
	}

	return nil
}
