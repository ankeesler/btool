package handlers_test

import (
	"reflect"
	"testing"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/pipeline"
	"github.com/ankeesler/btool/node/pipeline/handlers"
	"github.com/ankeesler/btool/node/testutil"
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
	ctx := pipeline.NewCtx()
	ctx.Nodes = []*node.Node{
		mainc,
		dep1h,
		dep0h,
	}
	if err := h.Handle(ctx); err != nil {
		t.Error(err)
	}

	ex := []*node.Node{
		dep0h,
		dep1h,
		mainc,
	}
	if !reflect.DeepEqual(ex, ctx.Nodes) {
		t.Error(ex, "!=", ctx.Nodes)
	}

	// Sad.
	dep0h.Dependencies = []*node.Node{mainc}
	if err := h.Handle(ctx); err == nil {
		t.Error("expected failure")
	}
}
