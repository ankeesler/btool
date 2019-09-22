// Package currenter provides a type that tells if a node.Node is up-to-date.
package currenter

import (
	"os"
	"time"

	"github.com/ankeesler/btool/log"
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
			log.Debugf("%s does not exist", n)
			return false, nil
		} else {
			return false, errors.Wrap(err, "lstat")
		}
	}

	latestT := time.Time{}
	var latestN *node.Node
	if err := node.Visit(n, func(vn *node.Node) error {
		vnInfo, err := os.Lstat(vn.Name)
		if err != nil {
			return errors.Wrap(err, "lstat")
		}

		modT := vnInfo.ModTime()
		if modT.After(latestT) {
			latestT = modT
			latestN = vn
		}

		return nil
	}); err != nil {
		return false, errors.Wrap(err, "visit")
	}

	if latestT.After(nInfo.ModTime()) {
		log.Debugf("%s has a newer dependency: %s", n, latestN)
		return false, nil
	}

	return true, nil
}
