package collector

import (
	"github.com/ankeesler/btool/node"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// Labels is a data structure that stores unmarshaled node.Node Label data.
// This data is common to all node.Node's.
type Labels struct {
	// Whether or not the node.Node is from a local project.
	Local bool `mapstructure:"io.btool.collector.local"`
}

// ToLabels will marshal a node.Node's Labels to the provided struct or return
// an error if it cannot do so. Extra fields from the node.Node's Labels are
// ignored.
func ToLabels(n *node.Node, i interface{}) error {
	return decode(i, &n.Labels)
}

// MustToLabels calls ToLabels and panics if it fails.
func MustToLabels(n *node.Node, i interface{}) {
	if err := ToLabels(n, i); err != nil {
		panic(err)
	}
}

// FromLabels will marshal a provided struct into the provided node.Node's
// Labels. This function will return an error if it fails to do so. The provided
// struct must be a pointer.
func FromLabels(n *node.Node, i interface{}) error {
	return decode(n.Labels, i)
}

// MustFromLabels calls FromLabels and panics if it fails.
func MustFromLabels(n *node.Node, i interface{}) {
	if err := FromLabels(n, i); err != nil {
		panic(err)
	}
}

func decode(from, to interface{}) error {
	if err := mapstructure.Decode(from, to); err != nil {
		return errors.Wrap(err, "decode")
	}
	return nil
}
