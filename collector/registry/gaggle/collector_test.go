package gaggle_test

import (
	"testing"

	"github.com/ankeesler/btool/collector/collectorfakes"
	"github.com/ankeesler/btool/collector/registry/gaggle"
	"github.com/ankeesler/btool/collector/testutil"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/nodefakes"
	"github.com/ankeesler/btool/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectorCollect(t *testing.T) {
	compilerCR := &nodefakes.FakeResolver{}
	linkerCR := &nodefakes.FakeResolver{}
	rf := &collectorfakes.FakeResolverFactory{}
	rf.NewCompileCReturnsOnCall(0, compilerCR)
	rf.NewLinkCReturnsOnCall(0, linkerCR)

	s := testutil.FakeStore()
	g := &registry.Gaggle{
		Nodes: []*registry.Node{
			&registry.Node{
				Name:         "a.h",
				Dependencies: []string{},
				Labels: map[string]interface{}{
					"io.btool.collector.cc.includePaths": []interface{}{
						"include/dir0",
						"include/dir1",
					},
					"io.btool.collector.cc.libraries": []interface{}{
						"some/library",
					},
				},
			},
			&registry.Node{
				Name:         "a.c",
				Dependencies: []string{"a.h"},
				Labels: map[string]interface{}{
					"io.btool.collector.cc.includePaths": []interface{}{
						"another/include/dir",
					},
				},
			},
			&registry.Node{
				Name:         "a.o",
				Dependencies: []string{"a.c"},
				Labels:       map[string]interface{}{},
				Resolver: registry.Resolver{
					Name: "compileC",
				},
			},
			&registry.Node{
				Name:         "a.a",
				Dependencies: []string{"a.o"},
				Labels:       map[string]interface{}{},
				Resolver: registry.Resolver{
					Name: "linkC",
				},
			},
		},
	}
	root := "/some/root"

	c := gaggle.New(rf)
	require.Nil(t, c.Collect(s, g, root))

	nodeAH := node.New("/some/root/a.h")
	nodeAH.Label(
		"io.btool.collector.cc.includePaths",
		[]string{
			"/some/root/include/dir0",
			"/some/root/include/dir1",
		},
	)
	nodeAH.Label(
		"io.btool.collector.cc.libraries",
		[]string{
			"/some/root/some/library",
		},
	)
	nodeAC := node.New("/some/root/a.c").Dependency(nodeAH)
	nodeAC.Label(
		"io.btool.collector.cc.includePaths",
		[]string{
			"/some/root/another/include/dir",
		},
	).Label(
		"io.btool.collector.cc.libraries",
		[]string{},
	)
	nodeAO := node.New("/some/root/a.o").Dependency(nodeAC)
	nodeAO.Label(
		"io.btool.collector.cc.includePaths",
		[]string{},
	).Label(
		"io.btool.collector.cc.libraries",
		[]string{},
	)
	nodeAO.Resolver = compilerCR
	nodeAA := node.New("/some/root/a.a").Dependency(nodeAO)
	nodeAA.Label(
		"io.btool.collector.cc.includePaths",
		[]string{},
	).Label(
		"io.btool.collector.cc.libraries",
		[]string{},
	)
	nodeAA.Resolver = linkerCR
	assert.Equal(t, nodeAH, s.Get(nodeAH.Name))
	assert.Equal(t, nodeAC, s.Get(nodeAC.Name))
	assert.Equal(t, nodeAO, s.Get(nodeAO.Name))
	assert.Equal(t, nodeAA, s.Get(nodeAA.Name))

	assert.Equal(
		t,
		[]string{
			"/some/root/include/dir0",
			"/some/root/include/dir1",
			"/some/root/another/include/dir",
		},
		rf.NewCompileCArgsForCall(0),
	)

	//assert.Equal(
	//	t,
	//	[]string{
	//		"/some/root/include/dir",
	//		"/some/root/another/include/dir",
	//	},
	//	ctx.IncludePaths(),
	//)

	//assert.Equal(
	//	t,
	//	[]*node.Node{
	//		nodeAA,
	//	},
	//	ctx.Libraries("a.h"),
	//)
}
