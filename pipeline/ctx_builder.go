package pipeline

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

// Root sets the root on the CtxBuilder's Ctx, and then returns the CtxBuilder.
func (cb *CtxBuilder) Root(root string) *CtxBuilder {
	cb.ctx.KV[CtxRoot] = root
	return cb
}

// Build returns the CtxBuilder's Ctx.
func (cb *CtxBuilder) Build() *Ctx {
	return cb.ctx
}
