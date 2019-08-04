package handlers_test

import (
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

func TestObject(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

	data := []struct {
		name  string
		nodes testutil.Nodes
		ext   string
	}{
		{
			name:  "BasicC",
			nodes: testutil.BasicNodesC.Copy(),
			ext:   ".c",
		},
		{
			name:  "BasicCC",
			nodes: testutil.BasicNodesCC.Copy(),
			ext:   ".cc",
		},
	}

	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			root := "/root"
			cache := "/cache"
			target := "dep-1/dep-1.o"
			compilerC := "cc"
			compilerCC := "c++"

			h := handlers.NewObject()
			ctx := pipeline.NewCtxBuilder().Nodes(
				datum.nodes,
			).Root(
				root,
			).Cache(
				cache,
			).Target(
				target,
			).CompilerC(
				compilerC,
			).CompilerCC(
				compilerCC,
			).Build()

			var compiler string
			if datum.ext == ".c" {
				compiler = compilerC
			} else if datum.ext == ".cc" {
				compiler = compilerCC
			} else {
				t.Fatalf("unknown compiler for extension %s", datum.ext)
			}

			source := "dep-1/dep-1" + datum.ext
			sourceN := node.Find(source, ctx.Nodes)
			if sourceN == nil {
				t.Fatal()
			}
			objectN := node.New(target).Dependency(sourceN)
			objectN.Resolver = resolvers.NewCompile(root, compiler, []string{root})
			exNodes := append(datum.nodes.Copy().Cast(), objectN)

			h.Handle(ctx)
			if ctx.Err != nil {
				t.Error(ctx.Err)
			}

			if diff := deep.Equal(exNodes, ctx.Nodes); diff != nil {
				t.Error(diff)
			}
		})
	}
}
