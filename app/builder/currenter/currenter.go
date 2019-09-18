// Package currenter provides a type that tells if a node.Node is up-to-date.
package currenter

import (
	"os"
	"time"

	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
)

// Currenter can tell whether a node.Node is up-to-date.
// It uses filesystem calls to figure out if a node.Node is older than its
// dependencies.
type Currenter struct {
}

// New creates a new Currenter.
func New() *Currenter {
	return &Currenter{}
}

func (c *Currenter) Current(n *node.Node) (bool, error) {
	nInfo, err := os.Lstat(n.Name)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, errors.Wrap(err, "lstat")
		}
	}

	latestT := time.Time{}
	for _, dN := range n.Dependencies {
		dInfo, err := os.Lstat(dN.Name)
		if err != nil {
			return false, errors.Wrap(err, "lstat")
		}

		modT := dInfo.ModTime()
		if modT.After(latestT) {
			latestT = modT
		}
	}

	if latestT.After(nInfo.ModTime()) {
		return false, nil
	}

	return true, nil
}
