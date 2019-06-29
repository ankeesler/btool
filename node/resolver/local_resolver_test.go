package resolver_test

import (
	"testing"

	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/resolver"
	"github.com/ankeesler/btool/node/testutil"
	"github.com/go-test/deep"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func TestLocalHandle(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

	data := []struct {
		name      string
		nodes     []*node.Node
		exNodes   []*node.Node
		exSuccess bool
	}{
		{
			name:      "Basic",
			nodes:     testutil.BasicNodesC.WithoutDependencies(),
			exNodes:   testutil.BasicNodesC,
			exSuccess: true,
		},
	}
	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			testutil.PopulateFS(datum.nodes, fs)

			l := resolver.NewLocal(fs, "/")

			acNodes, err := l.Handle(datum.nodes)
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

			if diff := deep.Equal(datum.exNodes, acNodes); diff != nil {
				t.Error(diff)
			}
		})
	}
}
