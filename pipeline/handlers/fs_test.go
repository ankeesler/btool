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

func TestHandle(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

	data := []struct {
		name      string
		exNodes   testutil.Nodes
		exSuccess bool
	}{
		{
			name:      "BasicC",
			exNodes:   testutil.BasicNodesC.WithoutDependencies(),
			exSuccess: true,
		},
		{
			name:      "BasicCC",
			exNodes:   testutil.BasicNodesCC.WithoutDependencies(),
			exSuccess: true,
		},
	}
	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			root := "/"

			fs := afero.NewMemMapFs()
			datum.exNodes.PopulateFS(root, fs)

			h := handlers.NewWalker(fs)

			ctx := pipeline.NewCtxBuilder().Root(root).Build()
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
