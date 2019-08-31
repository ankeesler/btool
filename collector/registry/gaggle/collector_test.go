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
	linkerCR := &nodefakes.FakeResolver{}
	rf := &collectorfakes.FakeResolverFactory{}
	rf.NewCompileCReturnsOnCall(0, compilerCR)
	rf.NewLinkCReturnsOnCall(0, linkerCR)

	ctx := collector.NewCtx(ns, rf)

	g := &registry.Gaggle{
		Metadata: map[string]interface{}{
			"includePaths": []string{
				"include/dir",
				"another/include/dir",
			},
			"libraries": map[string]string{
				"a.h": "a.a",
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
			&registry.Node{
				Name:         "a.a",
				Dependencies: []string{"a.o"},
				Resolver: registry.Resolver{
					Name: "linkC",
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
	nodeAA := node.New("/some/root/a.a").Dependency(nodeAO)
	nodeAA.Resolver = linkerCR
	assert.Equal(t, nodeAH, ns.Find(nodeAH.Name))
	assert.Equal(t, nodeAC, ns.Find(nodeAC.Name))
	assert.Equal(t, nodeAO, ns.Find(nodeAO.Name))
	assert.Equal(t, nodeAA, ns.Find(nodeAA.Name))

	assert.Equal(
		t,
		[]string{
			"/some/root/include/dir",
			"/some/root/another/include/dir",
		},
		rf.NewCompileCArgsForCall(0),
	)

	assert.Equal(
		t,
		[]string{
			"/some/root/include/dir",
			"/some/root/another/include/dir",
		},
		ctx.IncludePaths(),
	)

	assert.Equal(
		t,
		[]*node.Node{
			nodeAA,
		},
		ctx.Libraries("a.h"),
	)
}
