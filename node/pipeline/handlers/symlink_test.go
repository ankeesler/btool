package handlers_test

import (
	"testing"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/nodefakes"
	"github.com/ankeesler/btool/node/pipeline"
	"github.com/ankeesler/btool/node/pipeline/handlers"
	"github.com/ankeesler/btool/node/pipeline/handlers/handlersfakes"
	"github.com/ankeesler/btool/node/testutil"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/require"
)

func TestSymlink(t *testing.T) {
	symlinkR := &nodefakes.FakeResolver{}

	rf := &handlersfakes.FakeResolverFactory{}
	rf.NewSymlinkReturnsOnCall(0, symlinkR)

	to := "some/path/to/symlink"
	from := "dep-1/dep-1.o"

	h := handlers.NewSymlink(rf, to, from)

	ctx := pipeline.NewCtx()
	ctx.Nodes = testutil.BasicNodesCO.Copy()
	require.Nil(t, h.Handle(ctx))

	fromN := node.Find(from, ctx.Nodes)
	require.NotNil(t, fromN)
	exNodes := testutil.BasicNodesCO.Copy()
	toN := node.New(to).Dependency(fromN)
	toN.Resolver = symlinkR
	exNodes = append(exNodes, toN)
	require.Nil(t, deep.Equal(exNodes.Cast(), ctx.Nodes))
}
