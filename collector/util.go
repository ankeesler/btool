package collector

import (
	"fmt"
	"reflect"

	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
)

// CollectLabels will walk a node.Node graph and collect Labels into
// an array. The Labels are determined by the provided label key. The Label type
// must be a []string.
func CollectLabels(n *node.Node, label string) ([]string, error) {
	labels := make([]string, 0)

	if err := node.Visit(n, func(vn *node.Node) error {
		ls, ok := vn.Labels[label]
		if !ok {
			return nil
		}

		lsSlice, ok := ls.([]string)
		if !ok {
			return fmt.Errorf("expected []string, got %s", reflect.TypeOf(ls))
		}

		for _, l := range lsSlice {
			labels = append(labels, l)
		}

		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "visit")
	}

	return labels, nil
}

// AppendLabel will append a string to an existing node.Node's Label. The Label
// must be of type []string, otherwise an error will be returned.
func AppendLabel(n *node.Node, label string, apend string) error {
	l, ok := n.Labels[label]
	if !ok {
		l = make([]string, 0)
		n.Labels[label] = l
	}

	lSlice, ok := l.([]string)
	if !ok {
		return fmt.Errorf("expected []string, got %s", reflect.TypeOf(l))
	}

	lSlice = append(lSlice, apend)
	n.Labels[label] = lSlice

	return nil
}

// BoolLabel will get a Label from a node.Node, cast it to a bool, and return
// the value. If the Label does not exist on the node.Node, then this function
// will return false. If the Label does exist on the node.Node but it is not of
// type bool, then this function will return an error.
func BoolLabel(n *node.Node, label string) (bool, error) {
	if l, ok := n.Labels[label]; !ok {
		return false, nil
	} else if lBool, ok := l.(bool); !ok {
		return false, fmt.Errorf("expected bool, got %s", reflect.TypeOf(l))
	} else {
		return lBool, nil
	}
}
