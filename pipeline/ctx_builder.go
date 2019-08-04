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

// Root sets the root on the CtxBuilder's Ctx, and then returns the CtxBuilder.
func (cb *CtxBuilder) Root(root string) *CtxBuilder {
	cb.ctx.KV[CtxRoot] = root
	return cb
}

// Build returns the CtxBuilder's Ctx.
func (cb *CtxBuilder) Build() *Ctx {
	return cb.ctx
}
