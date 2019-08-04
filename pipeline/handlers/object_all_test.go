package handlers_test

import (
	"testing"

	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/node/testutil"
	"github.com/ankeesler/btool/pipeline"
	"github.com/ankeesler/btool/pipeline/handlers"
	"github.com/go-test/deep"
	"github.com/sirupsen/logrus"
)

func TestObjectAll(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

	data := []struct {
		name    string
		nodes   testutil.Nodes
		exNodes testutil.Nodes
	}{
		{
			name:    "BasicC",
			nodes:   testutil.BasicNodesC,
			exNodes: testutil.BasicNodesCO,
		},
		{
			name:    "BasicCC",
			nodes:   testutil.BasicNodesCC,
			exNodes: testutil.BasicNodesCCO,
		},
	}
	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			h := handlers.NewObjectAll()
			ctx := pipeline.NewCtxBuilder().Nodes(datum.nodes).Build()
			h.Handle(ctx)
			if ctx.Err != nil {
				t.Error(ctx.Err)
			}

			if diff := deep.Equal(datum.exNodes.Cast(), ctx.Nodes); diff != nil {
				t.Error(diff)
			}
		})
	}
}
