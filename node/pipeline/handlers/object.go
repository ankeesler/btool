package handlers

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/pipeline"
)

type object struct {
	s       Store
	rf      ResolverFactory
	project string
	target  string
}

// NewObject creates a pipeline.Handler that will add object file node.Node's to
// the node.Node list based on an object target.
func NewObject(
	s Store,
	rf ResolverFactory,
	project string,
	target string,
) pipeline.Handler {
	return &object{
		s:       s,
		rf:      rf,
		project: project,
		target:  target,
	}
}

func (o *object) Handle(ctx pipeline.Ctx) error {
	if !strings.HasSuffix(o.target, ".o") {
		return nil
	}

	var dN *node.Node
	sourceCN := node.Find(strings.ReplaceAll(o.target, ".o", ".c"), ctx.All())
	sourceCCN := node.Find(strings.ReplaceAll(o.target, ".o", ".cc"), ctx.All())
	if sourceCN != nil && sourceCCN != nil {
		return fmt.Errorf(
			"ambiguous object %s (%s or %s)",
			o.target,
			sourceCN.Name,
			sourceCCN.Name,
		)
	} else if sourceCN != nil {
		dN = sourceCN
	} else if sourceCCN != nil {
		dN = sourceCCN
	} else {
		return fmt.Errorf("unknown source for object %s", o.target)
	}

	objectN := objectNFromSourceN(o.s, o.rf, o.project, dN)
	ctx.Add(objectN)

	return nil
}

func (o *object) String() string { return "object" }

func objectNFromSourceN(
	s Store,
	rf ResolverFactory,
	project string,
	sourceN *node.Node,
) *node.Node {
	ext := filepath.Ext(sourceN.Name)

	includeDirs := []string{s.ProjectDir(project)}
	var r node.Resolver
	if ext == ".cc" {
		r = rf.NewCompileCC(includeDirs)
	} else {
		r = rf.NewCompileC(includeDirs)
	}

	object := strings.ReplaceAll(sourceN.Name, ext, ".o")
	log.Debugf(
		"adding %s -> %s",
		object,
		sourceN.Name,
	)
	objectN := node.New(object).Dependency(sourceN)
	objectN.Resolver = r
	return objectN
}
