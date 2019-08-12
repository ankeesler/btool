// Package resolvermapper provides a type that can map from a string to a
// node.Resolver.
package resolvermapper

import (
	"fmt"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/pipeline"
	"github.com/ankeesler/btool/node/resolvers"
)

// ResolverMapper can return a node.Resolver given a string and some config.
type ResolverMapper struct {
	ctx *pipeline.Ctx
}

// New creates a new ResolverMapper.
func New(ctx *pipeline.Ctx) *ResolverMapper {
	return &ResolverMapper{
		ctx: ctx,
	}
}

// Map returns a node.Resolver given an name and a config map. If no know
// Resolver can be created from the name and the config map, an error should
// be returned.
func (rm *ResolverMapper) Map(
	name string,
	config map[string]interface{},
) (node.Resolver, error) {
	root := rm.ctx.KV[pipeline.CtxRoot]
	compilerC := rm.ctx.KV[pipeline.CtxCompilerC]
	compilerCC := rm.ctx.KV[pipeline.CtxCompilerCC]
	archiver := rm.ctx.KV[pipeline.CtxArchiver]
	linker := rm.ctx.KV[pipeline.CtxLinker]
	switch name {
	case "compileC":
		return resolvers.NewCompile(root, compilerC, []string{root}), nil
	case "compilerCC":
		return resolvers.NewCompile(root, compilerCC, []string{root}), nil
	case "archive":
		return resolvers.NewArchive(root, archiver), nil
	case "link":
		return resolvers.NewLink(root, linker), nil
	default:
		return nil, fmt.Errorf("unknown resolver: %s", name)
	}
}
