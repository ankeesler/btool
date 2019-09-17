package collector

import (
	"strings"

	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
)

// CollectLabels will walk a node.Node graph and collect Labels into
// an array. The Labels are determined by the provided label key. The Label type
// must be an array.
func CollectLabels(n *node.Node, label string) ([]string, error) {
	labels := make([]string, 0)

	if err := node.Visit(n, func(vn *node.Node) error {
		if ls, ok := vn.Labels[label]; ok {
			for _, l := range strings.Split(ls, ",") {
				if l != "" {
					labels = append(labels, l)
				}
			}
		}
		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "visit")
	}

	return labels, nil
}
