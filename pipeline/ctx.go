package pipeline

import "fmt"

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

// Ctx provides a key-value store of information about a particular Pipeline
// instance. Members of the Pipeline can Get or Set keys on this store to provide
// information to other members of the Pipeline.
type Ctx struct {
	kv map[string][]string
}

// NewCtx creates a new Ctx with an empty key-value store.
func NewCtx() *Ctx {
	return &Ctx{
		kv: make(map[string][]string),
	}
}

// Get returns the value associated with the key in the Ctx instance. Iff nil is
// returned, then there is no value associated with this key.
func (ctx *Ctx) Get(key string) []string {
	return ctx.kv[key]
}

// Set appends the provided value to array associated with the provided key. The
// same Ctx object is returned in order to make calling this multiple times
// easier.
//   ctx := NewCtx().Append("a", "b").Append("a", "c").Append("z", "y")
func (ctx *Ctx) Append(key, value string) *Ctx {
	valueCurrent := ctx.kv[key]
	if valueCurrent == nil {
		valueCurrent = make([]string, 0)
		ctx.kv[key] = valueCurrent
	}
	valueCurrent = append(valueCurrent, value)
	return ctx
}

// String returns a string representation of the Ctx.
func (ctx *Ctx) String() string {
	return fmt.Sprintf("%s", ctx.kv)
}
