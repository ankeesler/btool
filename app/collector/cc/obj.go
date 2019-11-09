package cc

import (
	"path/filepath"

	"github.com/ankeesler/btool/app/collector"
	"github.com/ankeesler/btool/log"
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
	log.Debugf("collected include paths %s for node %s", includePaths, n)

	var r node.Resolver
	if ext == ".cc" {
		r = o.rf.NewCompileCC(includePaths)
	} else {
		r = o.rf.NewCompileC(includePaths)
	}

	on := node.New(replaceExt(n.Name, ext, ".o"))
	on.Dependency(n)
	on.Resolver = r

	s.Set(on)

	return nil
}

// CollectIncludePaths will walk a node.Node graph and return all of the include
// paths encoutered as a part of a node.Node's Labels along the way.
func CollectIncludePaths(n *node.Node) ([]string, error) {
	return collectLabels(n, func(l *Labels) []string { return l.IncludePaths })
}
