// Package unmarshaler provides a type that can unmarshal an API Node into a
// node.Node.
package unmarshaler

import (
	"fmt"
	"path/filepath"

	"github.com/ankeesler/btool/app/collector"
	"github.com/ankeesler/btool/node"
	nodev1 "github.com/ankeesler/btool/node/api/v1"
)

// Unmarshaler is a type that can unmarshal an API Node into a node.Node.
type Unmarshaler struct {
}

func New() *Unmarshaler {
	return &Unmarshaler{}
}

// Unmarshal will translate an API Node into a node.Node.
func (u *Unmarshaler) Unmarshal(
	s collector.Store,
	apiN *nodev1.Node,
	root string,
) (*node.Node, error) {
	n := node.New(filepath.Join(root, apiN.Name))
	for _, apiD := range apiN.Dependencies {
		d := s.Get(filepath.Join(root, apiD))
		if d == nil {
			return nil, fmt.Errorf("unknown dependencies %s for node %s", apiD, apiN)
		}

		n.Dependency(d)
	}

	return n, nil
}
