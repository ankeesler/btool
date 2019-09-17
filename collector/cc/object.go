package cc

import (
	"path/filepath"
	"strings"

	"github.com/ankeesler/btool/collector"
	"github.com/ankeesler/btool/node"
)

type Object struct {
	rf collector.ResolverFactory
}

func NewObject(rf collector.ResolverFactory) *Object {
	return &Object{
		rf: rf,
	}
}

func (o *Object) Consume(s collector.Store, n *node.Node) error {
	ext := filepath.Ext(n.Name)
	if ext != ".cc" && ext != ".c" {
		return nil
	}

	// TODO: is this bad to collect include paths from dependencies first?
	includePaths := make([]string, 0)
	node.Visit(n, func(vn *node.Node) error {
		if ips, ok := vn.Labels[LabelIncludePaths]; ok {
			// TODO: this is jank, we should have more of a better interface for this.
			for _, ip := range strings.Split(ips, ",") {
				includePaths = append(includePaths, ip)
			}
		}
		return nil
	})

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
