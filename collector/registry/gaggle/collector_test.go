package gaggle_test

import (
	"testing"

	"github.com/ankeesler/btool/collector"
	"github.com/ankeesler/btool/collector/collectorfakes"
	"github.com/ankeesler/btool/collector/registry/gaggle"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/nodefakes"
	"github.com/ankeesler/btool/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectorCollect(t *testing.T) {
	ns := collector.NewNodeStore(nil)

	compilerCR := &nodefakes.FakeResolver{}
	rf := &collectorfakes.FakeResolverFactory{}
	rf.NewCompileCReturnsOnCall(0, compilerCR)

	ctx := collector.NewCtx(ns, rf)

	g := &registry.Gaggle{
		Metadata: map[string]interface{}{
			"includePaths": map[string]string{
				"a.h": "include/dir",
				"b.h": "another/include/dir",
			},
		},
		Nodes: []*registry.Node{
			&registry.Node{
				Name:         "a.h",
				Dependencies: []string{},
			},
			&registry.Node{
				Name:         "a.c",
				Dependencies: []string{"a.h"},
			},
			&registry.Node{
				Name:         "a.o",
				Dependencies: []string{"a.c"},
				Resolver: registry.Resolver{
					Name: "compileC",
				},
			},
		},
	}
	root := "/some/root"
	c := gaggle.New()
	require.Nil(t, c.Collect(ctx, g, root))

	nodeAH := node.New("/some/root/a.h")
	nodeAC := node.New("/some/root/a.c").Dependency(nodeAH)
	nodeAO := node.New("/some/root/a.o").Dependency(nodeAC)
	nodeAO.Resolver = compilerCR
	assert.Equal(t, nodeAH, ns.Find(nodeAH.Name))
	assert.Equal(t, nodeAC, ns.Find(nodeAC.Name))
	assert.Equal(t, nodeAO, ns.Find(nodeAO.Name))

	assert.Equal(
		t,
		[]string{
			"include/dir",
		},
		rf.NewCompileCArgsForCall(0),
	)
}
