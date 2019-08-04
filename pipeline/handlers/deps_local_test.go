package handlers_test

import (
	"testing"

	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/node/testutil"
	"github.com/ankeesler/btool/pipeline"
	"github.com/ankeesler/btool/pipeline/handlers"
	"github.com/go-test/deep"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func TestDepsLocal(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

	data := []struct {
		name      string
		nodes     testutil.Nodes
		exNodes   testutil.Nodes
		exSuccess bool
	}{
		{
			name:      "Basic",
			nodes:     testutil.BasicNodesC.Copy().WithoutDependencies(),
			exNodes:   testutil.BasicNodesC.Copy(),
			exSuccess: true,
		},
	}
	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			root := "/"
			fs := afero.NewMemMapFs()
			datum.exNodes.PopulateFS(root, fs)

			h := handlers.NewDepsLocal(fs)

			ctx := pipeline.NewCtxBuilder().Nodes(datum.nodes).Root(root).Build()
			h.Handle(ctx)
			if datum.exSuccess {
				if ctx.Err != nil {
					t.Fatal(ctx.Err)
				}
			} else {
				if ctx.Err == nil {
					t.Fatal("expected failure")
				}
				return
			}

			if diff := deep.Equal(datum.exNodes.Cast(), ctx.Nodes); diff != nil {
				t.Error(diff)
			}
		})
	}
}
