package handlers_test

import (
	"path/filepath"
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

func TestObject(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

	data := []struct {
		name  string
		nodes testutil.Nodes
	}{
		{
			name:  "BasicC",
			nodes: testutil.BasicNodesC.Copy(),
		},
		{
			name:  "BasicCC",
			nodes: testutil.BasicNodesCC.Copy(),
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
			var ext string
			if strings.HasSuffix(datum.name, "CC") {
				compiler = compilerCC
				ext = ".cc"
			} else {
				compiler = compilerC
				ext = ".c"
			}

			source := "dep-1/dep-1" + ext
			sourceN := node.Find(source, ctx.Nodes)
			if sourceN == nil {
				t.Fatal()
			}
			name := filepath.Join(cache, filepath.Base(root), target)
			objectN := node.New(name).Dependency(sourceN)
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
