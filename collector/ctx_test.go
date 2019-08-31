package collector_test

import (
	"testing"

	"github.com/ankeesler/btool/collector"
	"github.com/ankeesler/btool/node"
	"github.com/stretchr/testify/assert"
)

func TestCtx(t *testing.T) {
	// a -> b, c
	// b -> c
	// c
	cN := node.New("c")
	bN := node.New("b").Dependency(cN)
	aN := node.New("a").Dependency(bN, cN)

	ns := collector.NewNodeStore(nil)
	ctx := collector.NewCtx(ns, nil)

	ctx.AddIncludePath("a-inc")
	ctx.AddIncludePath("b-inc")
	ctx.AddIncludePath("c-inc")
	assert.Equal(t, []string{"a-inc", "b-inc", "c-inc"}, ctx.IncludePaths())

	ctx.AddLibrary("a-inc", aN)
	ctx.AddLibrary("a-inc", bN)
	ctx.AddLibrary("c-inc", cN)
	assert.Equal(t, []*node.Node{aN, bN}, ctx.Libraries("a-inc"))
	assert.Equal(t, []*node.Node{cN}, ctx.Libraries("c-inc"))
	assert.Equal(t, []*node.Node(nil), ctx.Libraries("d-inc"))
}
