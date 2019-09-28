// Package cleaner provides a type that can remove all node.Node's from a
// filesystem.
package cleaner

import (
	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
)

// TODO: cleaner should not clean downloaded cache files.
// Makes it easier to run offline.

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
	return node.Visit(n, c.visit)
}

func (c *Cleaner) visit(n *node.Node) error {
	if n.Resolver != nil {
		c.cb.OnClean(n)
		if err := c.ra.RemoveAll(n.Name); err != nil {
			return errors.Wrap(err, "remove all")
		}
	}
	return nil
}
