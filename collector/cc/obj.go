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

func (o *Obj) Consume(s collector.Store, n *node.Node) error {
	ext := filepath.Ext(n.Name)
	if ext != ".cc" && ext != ".c" {
		return nil
	}

	// TODO: is this bad to collect include paths from dependencies first?
	includePaths, err := collector.CollectLabels(n, LabelIncludePaths)
	if err != nil {
		return errors.Wrap(err, "collect labels")
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
