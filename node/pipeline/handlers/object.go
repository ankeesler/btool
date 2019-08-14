package handlers

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/pipeline"
	"github.com/ankeesler/btool/node/resolvers"
	"github.com/sirupsen/logrus"
)

type object struct {
	s Store
}

// NewObject creates a pipeline.Handler that will add object file node.Node's to
// the node.Node list based on an object target.
func NewObject(s Store) pipeline.Handler {
	return &object{
		s: s,
	}
}

func (o *object) Handle(ctx *pipeline.Ctx) error {
	target := ctx.KV[pipeline.CtxTarget]
	if !strings.HasSuffix(target, ".o") {
		return nil
	}

	var dN *node.Node
	sourceCN := node.Find(strings.ReplaceAll(target, ".o", ".c"), ctx.Nodes)
	sourceCCN := node.Find(strings.ReplaceAll(target, ".o", ".cc"), ctx.Nodes)
	if sourceCN != nil && sourceCCN != nil {
		return fmt.Errorf(
			"ambiguous object %s (%s or %s)",
			target,
			sourceCN.Name,
			sourceCCN.Name,
		)
	} else if sourceCN != nil {
		dN = sourceCN
	} else if sourceCCN != nil {
		dN = sourceCCN
	} else {
		return fmt.Errorf("unknown source for object %s", target)
	}

	objectN := objectNFromSourceN(ctx, dN)
	ctx.Nodes = append(ctx.Nodes, objectN)

	symlinkN := node.New(target).Dependency(objectN)
	symlinkN.Resolver = resolvers.NewSymlink()
	ctx.Nodes = append(ctx.Nodes, symlinkN)

	return nil
}

func (o *object) Name() string { return "object" }

func objectNFromSourceN(ctx *pipeline.Ctx, sourceN *node.Node) *node.Node {
	root := ctx.KV[pipeline.CtxRoot]
	compiler := getCompiler(ctx, sourceN)

	ext := filepath.Ext(sourceN.Name)
	object := cacheObjectPath(ctx, strings.ReplaceAll(sourceN.Name, ext, ".o"))

	logrus.Debugf(
		"adding %s -> %s with compiler %s",
		object,
		sourceN.Name,
		compiler,
	)
	objectN := node.New(object).Dependency(sourceN)
	objectN.Resolver = resolvers.NewCompile(root, compiler, []string{root})

	return objectN
}

func getCompiler(ctx *pipeline.Ctx, n *node.Node) string {
	ext := filepath.Ext(n.Name)
	switch ext {
	case ".c":
		return ctx.KV[pipeline.CtxCompilerC]
	case ".cc":
		return ctx.KV[pipeline.CtxCompilerCC]
	default:
		return ""
	}
}
