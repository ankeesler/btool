package pipeline

import "github.com/ankeesler/btool/node"

// CtxBuilder is a utility type used to make creating Ctx's easier.
type CtxBuilder struct {
	ctx *Ctx
}

// NewCtxBuilder creates a new CtxBuilder with a default Ctx.
func NewCtxBuilder() *CtxBuilder {
	return &CtxBuilder{
		ctx: NewCtx(),
	}
}

// Nodes sets the Nodes on the CtxBuilder's Ctx, and then returns the CtxBuilder.
func (cb *CtxBuilder) Nodes(nodes []*node.Node) *CtxBuilder {
	cb.ctx.Nodes = nodes
	return cb
}

// Project sets the root on the CtxBuilder's Ctx, and then returns the CtxBuilder.
func (cb *CtxBuilder) Project(project string) *CtxBuilder {
	cb.ctx.KV[CtxProject] = project
	return cb
}

// Root sets the root on the CtxBuilder's Ctx, and then returns the CtxBuilder.
func (cb *CtxBuilder) Root(root string) *CtxBuilder {
	cb.ctx.KV[CtxRoot] = root
	return cb
}

// Cache sets the cache on the CtxBuilder's Ctx, and then returns the CtxBuilder.
func (cb *CtxBuilder) Cache(cache string) *CtxBuilder {
	cb.ctx.KV[CtxCache] = cache
	return cb
}

// Target sets the target on the CtxBuilder's Ctx, and then returns the CtxBuilder.
func (cb *CtxBuilder) Target(target string) *CtxBuilder {
	cb.ctx.KV[CtxTarget] = target
	return cb
}

// CompileC sets the C compiler on the CtxBuilder's Ctx, and then returns the
// CtxBuilder.
func (cb *CtxBuilder) CompilerC(compiler string) *CtxBuilder {
	cb.ctx.KV[CtxCompilerC] = compiler
	return cb
}

// CompileC sets the C++ compiler on the CtxBuilder's Ctx, and then returns the
// CtxBuilder.
func (cb *CtxBuilder) CompilerCC(compiler string) *CtxBuilder {
	cb.ctx.KV[CtxCompilerCC] = compiler
	return cb
}

// Archiver sets the archiver on the CtxBuilder's Ctx, and then returns the CtxBuilder.
func (cb *CtxBuilder) Archiver(archiver string) *CtxBuilder {
	cb.ctx.KV[CtxArchiver] = archiver
	return cb
}

// Linker sets the linker on the CtxBuilder's Ctx, and then returns the CtxBuilder.
func (cb *CtxBuilder) Linker(linker string) *CtxBuilder {
	cb.ctx.KV[CtxLinker] = linker
	return cb
}

// Build returns the CtxBuilder's Ctx.
func (cb *CtxBuilder) Build() *Ctx {
	return cb.ctx
}
