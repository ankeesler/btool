package pipeline

import (
	"fmt"

	"github.com/ankeesler/btool/node"
)

// The following are core pieces of information that apply to every member of
// the Pipeline.
const (
	CtxRoot   = "pipeline.root"
	CtxCache  = "pipeline.cache"
	CtxTarget = "pipeline.target"

	CtxCompilerC  = "pipeline.compiler.c"
	CtxCompilerCC = "pipeline.compiler.cc"
	CtxLinker     = "pipeline.linker"
)

// Ctx provides 3 things:
//   - the node.Node list on which this Pipeline is operating
//   - an error that represents whether there has been a failure
//   - a key-value store of information about a particular Pipeline
type Ctx struct {
	Nodes []*node.Node
	Err   error
	KV    map[string]string
}

// NewCtx creates a new Ctx with an empty node.Node list, nil error, and empty
// key-value store.
func NewCtx() *Ctx {
	return &Ctx{
		Nodes: make([]*node.Node, 0),
		Err:   nil,
		KV:    make(map[string]string),
	}
}

// String returns a string representation of the Ctx.
func (ctx *Ctx) String() string {
	return fmt.Sprintf("%s/%s/%s", ctx.Nodes, ctx.Err, ctx.KV)
}
