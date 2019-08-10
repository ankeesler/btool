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

func TestFS(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

	data := []struct {
		name      string
		exNodes   testutil.Nodes
		exSuccess bool
	}{
		{
			name:      "BasicC",
			exNodes:   testutil.BasicNodesC.Copy().WithoutDependencies(),
			exSuccess: true,
		},
		{
			name:      "BasicCC",
			exNodes:   testutil.BasicNodesCC.Copy().WithoutDependencies(),
			exSuccess: true,
		},
	}
	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			root := "/"

			fs := afero.NewMemMapFs()
			datum.exNodes.PopulateFS(root, fs)

			h := handlers.NewFS(fs)

			ctx := pipeline.NewCtxBuilder().Root(root).Build()
			err := h.Handle(ctx)
			if datum.exSuccess {
				if err != nil {
					t.Fatal(err)
				}
			} else {
				if err == nil {
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
