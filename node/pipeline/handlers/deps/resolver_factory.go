package deps

import (
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/pipeline"
	"github.com/ankeesler/btool/node/resolvers"
)

type resolverFactory struct {
	ctx *pipeline.Ctx
}

func newResolverFactory(ctx *pipeline.Ctx) *resolverFactory {
	return &resolverFactory{
		ctx: ctx,
	}
}

func (rf *resolverFactory) make(resolver string) node.Resolver {
	root := rf.ctx.KV[pipeline.CtxRoot]
	compilerC := rf.ctx.KV[pipeline.CtxCompilerC]
	compilerCC := rf.ctx.KV[pipeline.CtxCompilerCC]
	archiver := rf.ctx.KV[pipeline.CtxArchiver]
	linker := rf.ctx.KV[pipeline.CtxLinker]
	switch resolver {
	case "compileC":
		return resolvers.NewCompile(root, compilerC, []string{root})
	case "compilerCC":
		return resolvers.NewCompile(root, compilerCC, []string{root})
	case "archive":
		return resolvers.NewArchive(root, archiver)
	case "link":
		return resolvers.NewLink(root, linker)
	default:
		return nil
	}
}
