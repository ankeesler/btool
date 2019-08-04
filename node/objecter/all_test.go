package objecter_test

import (
	"testing"

	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/objecter"
	"github.com/ankeesler/btool/node/testutil"
	"github.com/go-test/deep"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func TestHandle(t *testing.T) {
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
			fs := afero.NewMemMapFs()
			datum.exNodes.PopulateFS("/", fs)

			a := objecter.NewAll()

			cfg := node.Config{}
			acNodes, err := a.Handle(&cfg, datum.nodes.Copy())
			if err != nil {
				t.Fatal(err)
			}

			if diff := deep.Equal(datum.exNodes.Cast(), acNodes); diff != nil {
				t.Error(diff)
			}
		})
	}
}
