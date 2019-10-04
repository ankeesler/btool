// Package unmarshaler provides a type that can unmarshal an API Node into a
// node.Node.
package unmarshaler

import (
	"fmt"
	"path/filepath"

	"github.com/ankeesler/btool/app/collector"
	"github.com/ankeesler/btool/app/collector/cc"
	"github.com/ankeesler/btool/node"
	nodev1 "github.com/ankeesler/btool/node/api/v1"
	"github.com/pkg/errors"
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

	if err := prependRoot(n, root); err != nil {
		return nil, errors.Wrap(err, "prepend root")
	}

	for _, apiD := range apiN.Dependencies {
		d := s.Get(filepath.Join(root, apiD))
		if d == nil {
			return nil, fmt.Errorf("unknown dependencies %s for node %s", apiD, apiN)
		}

		n.Dependency(d)
	}

	return n, nil
}

func prependRoot(n *node.Node, root string) error {
	// TODO: is this bad to reach across sibling packages?
	var labels cc.Labels
	if err := collector.FromLabels(n, &labels); err != nil {
		return errors.Wrap(err, "from labels")
	}

	if labels.IncludePaths == nil {
		labels.IncludePaths = []string{}
	}

	for i := range labels.IncludePaths {
		labels.IncludePaths[i] = filepath.Join(root, labels.IncludePaths[i])
	}

	if labels.Libraries == nil {
		labels.Libraries = []string{}
	}

	for i := range labels.Libraries {
		labels.Libraries[i] = filepath.Join(root, labels.Libraries[i])
	}

	if err := collector.ToLabels(n, &labels); err != nil {
		return errors.Wrap(err, "to labels")
	}

	return nil
}
