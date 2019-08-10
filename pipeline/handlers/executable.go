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

type executable struct {
}

// NewExecutable creates a pipeline.Handler that will add object file node.Node's
// to the node.Node list based on an executable target.
func NewExecutable() pipeline.Handler {
	return &executable{}
}

func (e *executable) Handle(ctx *pipeline.Ctx) error {
	target := ctx.KV[pipeline.CtxTarget]
	if filepath.Ext(target) != "" {
		return nil
	}

	var dN *node.Node
	sourceCN := node.Find(target+".c", ctx.Nodes)
	sourceCCN := node.Find(target+".cc", ctx.Nodes)
	if sourceCN != nil && sourceCCN != nil {
		return fmt.Errorf(
			"ambiguous executable %s (%s or %s)",
			target,
			sourceCN.Name,
			sourceCCN.Name,
		)
	} else if sourceCN != nil {
		dN = sourceCN
	} else if sourceCCN != nil {
		dN = sourceCCN
	} else {
		return fmt.Errorf("unknown source for executable %s", target)
	}

	objectNodes := make([]*node.Node, 0)
	objectNodes = collectObjects(ctx, dN, objectNodes)

	linker := ctx.KV[pipeline.CtxLinker]
	name := cacheExecutablePath(ctx, target)
	targetN := node.New(name)
	for _, objectN := range objectNodes {
		ctx.Nodes = append(ctx.Nodes, objectN)
		targetN.Dependency(objectN)
	}
	targetN.Resolver = resolvers.NewLink(ctx.KV[pipeline.CtxRoot], linker)
	ctx.Nodes = append(ctx.Nodes, targetN)

	symlinkN := node.New(target).Dependency(targetN)
	symlinkN.Resolver = resolvers.NewSymlink()
	ctx.Nodes = append(ctx.Nodes, symlinkN)

	return nil
}

func (e *executable) Name() string { return "executable" }

func collectObjects(
	ctx *pipeline.Ctx,
	sourceN *node.Node,
	objectNodes []*node.Node,
) []*node.Node {
	logrus.Debugf("collect objects from %s", sourceN.Name)

	objectN := objectNFromSourceN(ctx, sourceN)
	if node.Find(objectN.Name, objectNodes) != nil {
		return objectNodes
	}
	objectNodes = append(objectNodes, objectN)

	for _, dN := range sourceN.Dependencies {
		ext := filepath.Ext(sourceN.Name)
		source := strings.ReplaceAll(dN.Name, ".h", ext)
		if source == sourceN.Name {
			continue
		}

		sourceN := node.Find(source, ctx.Nodes)
		logrus.Debugf("dependency %s, source %s, found %s", dN, source, sourceN)
		if sourceN != nil {
			objectNodes = collectObjects(ctx, sourceN, objectNodes)
		}
	}

	return objectNodes
}
