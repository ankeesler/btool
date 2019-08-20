package handlers_test

import (
	"testing"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/pipeline/handlers"
	pipelinetestutil "github.com/ankeesler/btool/node/pipeline/testutil"
	"github.com/ankeesler/btool/node/testutil"
	"github.com/stretchr/testify/require"
)

func TestSortTopo(t *testing.T) {
	h := handlers.NewSortTopo()

	nodes := testutil.BasicNodesC.Copy()
	dep0h := node.Find("dep-0/dep-0.h", nodes)
	if dep0h == nil {
		t.Fatal()
	}
	dep1h := node.Find("dep-1/dep-1.h", nodes)
	if dep1h == nil {
		t.Fatal()
	}
	mainc := node.Find("main.c", nodes)
	if mainc == nil {
		t.Fatal()
	}

	// Happy.
	ctx := pipelinetestutil.NewCtx()
	ctx.Nodes = []*node.Node{
		mainc,
		dep1h,
		dep0h,
	}
	require.Nil(t, h.Handle(ctx))

	//ex := []*node.Node{
	//	dep0h,
	//	dep1h,
	//	mainc,
	//}
	//require.Equal(t, ex, ctx.Nodes)

	// Sad.
	dep0h.Dependencies = []*node.Node{mainc}
	require.NotNil(t, h.Handle(ctx))
}
