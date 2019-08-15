package handlers_test

import (
	"testing"

	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/node/pipeline"
	"github.com/ankeesler/btool/node/pipeline/handlers"
	"github.com/ankeesler/btool/node/testutil"
	"github.com/go-test/deep"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			exNodes:   testutil.BasicNodesC.Copy(),
			exSuccess: true,
		},
		{
			name:      "BasicCC",
			exNodes:   testutil.BasicNodesCC.Copy(),
			exSuccess: true,
		},
	}
	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			datum.exNodes.PopulateFS("/", fs)

			h := handlers.NewFS(fs, "/")

			ctx := pipeline.NewCtx()
			err := h.Handle(ctx)
			if datum.exSuccess {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err, "expected error to have occurred")
				return
			}

			assert.Nil(t, deep.Equal(datum.exNodes.Cast(), ctx.Nodes))
		})
	}
}
