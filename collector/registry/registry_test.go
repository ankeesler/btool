package registry_test

import (
	"testing"

	"github.com/ankeesler/btool/collector"
	"github.com/ankeesler/btool/collector/collectorfakes"
	"github.com/ankeesler/btool/collector/registry"
	"github.com/ankeesler/btool/collector/registry/registryfakes"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/nodefakes"
	registrypkg "github.com/ankeesler/btool/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistryCollect(t *testing.T) {
	g := &registryfakes.FakeGaggler{}
	g.GaggleReturnsOnCall(
		0,
		&registrypkg.Gaggle{
			Metadata: map[string]interface{}{
				"includeDirs": []string{
					"include/dir",
					"another/include/dir",
				},
			},
			Nodes: []*registrypkg.Node{
				&registrypkg.Node{
					Name:         "a.h",
					Dependencies: []string{},
				},
				&registrypkg.Node{
					Name:         "a.c",
					Dependencies: []string{"a.h"},
				},
				&registrypkg.Node{
					Name:         "a.o",
					Dependencies: []string{"a.c"},
					Resolver: registrypkg.Resolver{
						Name: "compileC",
					},
				},
			},
		},
	)
	g.RootReturns("/some/root")

	r := registry.New(g)

	ns := collector.NewNodeStore(nil)

	compilerCR := &nodefakes.FakeResolver{}
	rf := &collectorfakes.FakeResolverFactory{}
	rf.NewCompileCReturnsOnCall(0, compilerCR)

	ctx := collector.NewCtx(ns, rf)
	acNode := node.New("main")
	require.Nil(t, r.Collect(ctx, acNode))

	exNode := node.New("main")
	assert.Equal(t, exNode, acNode)

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
			"another/include/dir",
		},
		rf.NewCompileCArgsForCall(0),
	)
}
