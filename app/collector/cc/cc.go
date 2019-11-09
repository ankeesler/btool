// Package cc provides collector.Producer's and collector.Consumer's specific to
// C/C++ code.
package cc

import (
	"strings"

	"github.com/ankeesler/btool/app/collector"
	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
)

// Labels is a data structure that stores unmarshaled node.Node Label data.
// This data is common to all C/C++ related node.Node's.
type Labels struct {
	// Include paths required when compiling this thing.
	IncludePaths []string `mapstructure:"io.btool.collector.cc.includePaths"`

	// Libraries required when linking this thing.
	Libraries []string `mapstructure:"io.btool.collector.cc.libraries"`

	// Link flags added to linker invocations.
	LinkFlags []string `mapstructure:"io.btool.collector.cc.linkFlags"`
}

func collectLabels(
	n *node.Node,
	lFunc func(*Labels) []string,
) ([]string, error) {
	l := make([]string, 0)
	if err := node.Visit(n, func(vn *node.Node) error {
		var labels Labels
		if err := collector.FromLabels(vn, &labels); err != nil {
			return errors.Wrap(err, "from labels")
		}

		l = append(l, lFunc(&labels)...)

		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "visit")
	}

	return l, nil
}

func replaceExt(s, old, new string) string {
	i := strings.LastIndex(s, old)
	if i == -1 {
		return s
	}

	return s[:i] + new
}
