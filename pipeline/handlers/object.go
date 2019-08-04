package handlers

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/resolvers"
	"github.com/ankeesler/btool/pipeline"
	"github.com/sirupsen/logrus"
)

type object struct {
}

// NewObject creates a pipeline.Handler that will add object file node.Node's to
// the node.Node list based on an object target.
func NewObject() pipeline.Handler {
	return &object{}
}

func (o *object) Handle(ctx *pipeline.Ctx) {
	target := ctx.KV[pipeline.CtxTarget]
	if !strings.HasSuffix(target, ".o") {
		return
	}

	var dN *node.Node
	sourceCN := node.Find(strings.ReplaceAll(target, ".o", ".c"), ctx.Nodes)
	sourceCCN := node.Find(strings.ReplaceAll(target, ".o", ".cc"), ctx.Nodes)
	if sourceCN != nil && sourceCCN != nil {
		ctx.Err = fmt.Errorf(
			"ambiguous object %s (%s or %s)",
			target,
			sourceCN.Name,
			sourceCCN.Name,
		)
		return
	} else if sourceCN != nil {
		dN = sourceCN
	} else if sourceCCN != nil {
		dN = sourceCCN
	} else {
		ctx.Err = fmt.Errorf("unknown source for object %s", target)
		return
	}

	objectN := objectNFromSourceN(ctx, dN)
	ctx.Nodes = append(ctx.Nodes, objectN)
}

func (o *object) Name() string { return "object" }

func objectNFromSourceN(ctx *pipeline.Ctx, sourceN *node.Node) *node.Node {
	root := ctx.KV[pipeline.CtxRoot]
	compiler := getCompiler(ctx, sourceN)

	ext := filepath.Ext(sourceN.Name)
	object := strings.ReplaceAll(sourceN.Name, ext, ".o")

	logrus.Debugf("adding %s -> %s", object, sourceN.Name)
	objectN := node.New(object).Dependency(sourceN)
	objectN.Resolver = resolvers.NewCompile(root, compiler, []string{root})

	return objectN
}

func getCompiler(ctx *pipeline.Ctx, n *node.Node) string {
	ext := filepath.Ext(n.Name)
	switch ext {
	case ".c":
		return ctx.KV[pipeline.CtxCompilerCC]
	case ".cc":
		return ctx.KV[pipeline.CtxCompilerC]
	default:
		return ""
	}
}
