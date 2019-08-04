package handlers_test

import (
	"reflect"
	"testing"

	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/testutil"
	"github.com/ankeesler/btool/pipeline"
	"github.com/ankeesler/btool/pipeline/handlers"
	"github.com/sirupsen/logrus"
)

func TestSortTopo(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

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

	//dep0h := *testutil.Dep0h
	//dep1h := *testutil.Dep1h
	//mainc := *testutil.Mainc

	// Happy.
	ctx := pipeline.NewCtxBuilder().Nodes(
		[]*node.Node{
			mainc,
			dep1h,
			dep0h,
		}).Build()
	h.Handle(ctx)
	if ctx.Err != nil {
		t.Error(ctx.Err)
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
	h.Handle(ctx)
	if ctx.Err == nil {
		t.Error("expected failure")
	}
}
