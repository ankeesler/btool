package handlers_test

import (
	"strings"
	"testing"

	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/nodefakes"
	"github.com/ankeesler/btool/node/pipeline"
	"github.com/ankeesler/btool/node/pipeline/handlers"
	"github.com/ankeesler/btool/node/pipeline/handlers/handlersfakes"
	"github.com/ankeesler/btool/node/testutil"
	"github.com/go-test/deep"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecutable(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

	data := []struct {
		name             string
		nodes            testutil.Nodes
		nodesWithObjects testutil.Nodes
	}{
		{
			name:             "BasicC",
			nodes:            testutil.BasicNodesC.Copy(),
			nodesWithObjects: testutil.BasicNodesCO.Copy(),
		},
		{
			name:             "BasicCC",
			nodes:            testutil.BasicNodesCC.Copy(),
			nodesWithObjects: testutil.BasicNodesCCO.Copy(),
		},
	}
	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			s := &handlersfakes.FakeStore{}
			s.ProjectDirReturns("/some-project-dir")

			compileCR := &nodefakes.FakeResolver{}
			compileCCR := &nodefakes.FakeResolver{}
			linkR := &nodefakes.FakeResolver{}
			symlinkR := &nodefakes.FakeResolver{}

			rf := &handlersfakes.FakeResolverFactory{}
			rf.NewCompileCReturns(compileCR)
			rf.NewCompileCCReturns(compileCCR)
			rf.NewLinkReturnsOnCall(0, linkR)
			rf.NewSymlinkReturnsOnCall(0, symlinkR)

			h := handlers.NewExecutable(s, rf, "some-project", "main")
			ctx := pipeline.NewCtx()
			ctx.Nodes = datum.nodes
			require.Nil(t, h.Handle(ctx))

			name := "main"
			executableN := node.New(name)
			executableN.Resolver = linkR

			exNodes := datum.nodesWithObjects
			exNodes = append(exNodes, executableN)
			for _, n := range exNodes {
				if strings.HasSuffix(n.Name, ".o") {
					executableN.Dependency(n)
					if strings.HasSuffix(datum.name, "CC") {
						n.Resolver = compileCCR
					} else {
						n.Resolver = compileCR
					}
				}
			}

			symlinkN := node.New("main").Dependency(executableN)
			symlinkN.Resolver = symlinkR
			exNodes = append(exNodes, symlinkN)

			node.SortAlpha(ctx.Nodes)
			node.SortAlpha(exNodes)

			t.Log(ctx.Nodes)
			t.Log(exNodes)

			assert.Nil(t, deep.Equal(exNodes.Cast(), ctx.Nodes))
		})
	}
}
