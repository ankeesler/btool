// Package cleaner provides a type that can remove all node.Node's from a
// filesystem.
package cleaner

import (
	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . RemoveAller

// RemoveAller is a type that can delete a file from the filesystem. If the file
// does not exist, then it shouldn't fail.
type RemoveAller interface {
	RemoveAll(string) error
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Callback

// Callback is an interface that clients can use to be notified when a node.Node
// is removed from the filesystem.
type Callback interface {
	OnClean(*node.Node)
}

// Cleaner can remove all node.Node's from a filesystem.
type Cleaner struct {
	ra RemoveAller
	cb Callback
}

// New creates a new Cleaner.
func New(ra RemoveAller, cb Callback) *Cleaner {
	return &Cleaner{
		ra: ra,
		cb: cb,
	}
}

// Clean will walk a node.Node graph and remove all node.Node's from the
// filesystem.
func (c *Cleaner) Clean(n *node.Node) error {
	return c.clean(n, make(map[*node.Node]bool))
}

func (c *Cleaner) clean(n *node.Node, cleaned map[*node.Node]bool) error {
	if cleaned[n] {
		return nil
	}

	log.Debugf("cleaning %s", n.Name)

	for _, dN := range n.Dependencies {
		if err := c.clean(dN, cleaned); err != nil {
			return errors.Wrap(err, "clean")
		}
	}

	if n.Resolver != nil {
		c.cb.OnClean(n)
		if err := c.ra.RemoveAll(n.Name); err != nil {
			return errors.Wrap(err, "remove all")
		}
	}

	cleaned[n] = true

	return nil
}
