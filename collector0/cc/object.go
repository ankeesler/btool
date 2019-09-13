package cc

import (
	"path/filepath"
	"strings"

	collector "github.com/ankeesler/btool/collector0"
	"github.com/ankeesler/btool/node"
)

type Object struct {
	rf ResolverFactory
}

func NewObject(rf ResolverFactory) *Object {
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
		ips := vn.Labels[LabelIncludePaths]
		// TODO: this is jank, we should have more of a better interface for this.
		for _, ip := range strings.Split(ips, ",") {
			includePaths = append(includePaths, ip)
		}
		return nil
	})

	var r node.Resolver
	if ext == ".cc" {
		r = o.rf.NewCompileCC(includePaths)
	} else {
		r = o.rf.NewCompileC(includePaths)
	}
	n.Resolver = r

	on := node.New(strings.ReplaceAll(n.Name, ext, ".o"))
	on.Dependency(n)
	on.Resolver = r

	s.Set(on)

	return nil
}
