package resolvers

import (
	"fmt"
	"os"

	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
)

type symlink struct {
	// Note! afero does not contain support for symlinking currently.
	// See https://github.com/spf13/afero/pull/212/files.
	//fs afero.Fs
}

// NewSymlink returns a node.Resolver that symlinks a node.Node to its
// dependency.
func NewSymlink( /*fs afero.Fs*/ ) node.Resolver {
	return &symlink{
		//fs: fs,
	}
}

func (s *symlink) Resolve(n *node.Node) error {
	if len(n.Dependencies) != 1 {
		return fmt.Errorf("expected 1 dependency, got %d", len(n.Dependencies))
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
