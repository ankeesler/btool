package cc

import (
	"path/filepath"
	"strings"

	"github.com/ankeesler/btool/collector"
	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
)

// Obj is a type that can add .o node.Node's to a collector.Store given .c/.cc
// node.Node's.
type Obj struct {
	rf collector.ResolverFactory
}

// NewObj creates a new Obj.
func NewObj(rf collector.ResolverFactory) *Obj {
	return &Obj{
		rf: rf,
	}
}

// Consume will listen for .c/.cc files and create objects for them. It will walk
// the node.Node.Dependencies in order to collect a list of include paths for
// compilation.
func (o *Obj) Consume(s collector.Store, n *node.Node) error {
	ext := filepath.Ext(n.Name)
	if ext != ".cc" && ext != ".c" {
		return nil
	}

	includePaths, err := CollectIncludePaths(n)
	if err != nil {
		return errors.Wrap(err, "collect include paths")
	}

	var r node.Resolver
	if ext == ".cc" {
		r = o.rf.NewCompileCC(includePaths)
	} else {
		r = o.rf.NewCompileC(includePaths)
	}

	on := node.New(strings.ReplaceAll(n.Name, ext, ".o"))
	on.Dependency(n)
	on.Resolver = r

	s.Set(on)

	return nil
}

// CollectIncludePaths will walk a node.Node graph and return all of the include
// paths encoutered as a part of a node.Node's Labels along the way.
func CollectIncludePaths(n *node.Node) ([]string, error) {
	includePaths := make([]string, 0)
	if err := node.Visit(n, func(vn *node.Node) error {
		var labels Labels
		if err := collector.FromLabels(vn, &labels); err != nil {
			return errors.Wrap(err, "from labels")
		}

		includePaths = append(includePaths, labels.IncludePaths...)

		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "visit")
	}

	return includePaths, nil
}
