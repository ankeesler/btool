package walker_test

import (
	"testing"

	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/testutil"
	"github.com/ankeesler/btool/node/walker"
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
			fs := afero.NewMemMapFs()
			datum.exNodes.PopulateFS("/", fs)

			w := walker.New(fs)

			cfg := node.Config{
				Root: "/",
			}
			acNodes, err := w.Handle(&cfg, []*node.Node{})
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

			if diff := deep.Equal(datum.exNodes.Cast(), acNodes); diff != nil {
				t.Error(diff)
			}
		})
	}
}
