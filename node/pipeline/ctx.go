package pipeline

import (
	"fmt"

	"github.com/ankeesler/btool/node"
)

// Ctx provides 2 things:
//   - the node.Node list on which this Pipeline is operating
//   - a key-value store of information about a particular Pipeline
type Ctx struct {
	Nodes []*node.Node
	KV    map[string]string
}

// NewCtx creates a new Ctx with an empty node.Node list and an empty
// key-value store.
func NewCtx() *Ctx {
	return &Ctx{
		Nodes: make([]*node.Node, 0),
		KV:    make(map[string]string),
	}
}

// String returns a string representation of the Ctx.
func (ctx *Ctx) String() string {
	return fmt.Sprintf("%s: %s", ctx.Nodes, ctx.KV)
}
