package handlers_test

import (
	"strings"
	"testing"

	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/resolvers"
	"github.com/ankeesler/btool/node/testutil"
	"github.com/ankeesler/btool/pipeline"
	"github.com/ankeesler/btool/pipeline/handlers"
	"github.com/go-test/deep"
	"github.com/sirupsen/logrus"
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
			root := "/root"
			target := "main"
			compilerC := "cc"
			compilerCC := "c++"
			linker := "ld"

			h := handlers.NewExecutable()
			ctx := pipeline.NewCtxBuilder().Nodes(
				datum.nodes,
			).Root(
				root,
			).Target(
				target,
			).CompilerC(
				compilerC,
			).CompilerCC(
				compilerCC,
			).Linker(
				linker,
			).Build()
			h.Handle(ctx)
			if ctx.Err != nil {
				t.Error(ctx.Err)
			}

			var compiler string
			if strings.HasSuffix(datum.name, "CC") {
				compiler = compilerCC
			} else {
				compiler = compilerC
			}

			executableN := node.New(target)
			executableN.Resolver = resolvers.NewLink(root, linker)

			exNodes := datum.nodesWithObjects
			exNodes = append(exNodes, executableN)
			for _, n := range exNodes {
				if strings.HasSuffix(n.Name, ".o") {
					executableN.Dependency(n)
					n.Resolver = resolvers.NewCompile(root, compiler, []string{root})
				}
			}

			node.SortAlpha(ctx.Nodes)
			node.SortAlpha(exNodes)

			if diff := deep.Equal(exNodes.Cast(), ctx.Nodes); diff != nil {
				t.Error(diff)
			}
		})
	}
}
