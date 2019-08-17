package handlers_test

import (
	"strings"
	"testing"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/nodefakes"
	"github.com/ankeesler/btool/node/pipeline"
	"github.com/ankeesler/btool/node/pipeline/handlers"
	"github.com/ankeesler/btool/node/pipeline/handlers/handlersfakes"
	"github.com/ankeesler/btool/node/testutil"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestObject(t *testing.T) {
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
			s := &handlersfakes.FakeStore{}
			s.ProjectDirReturns("/some-project-dir")

			compileCR := &nodefakes.FakeResolver{}
			compileCCR := &nodefakes.FakeResolver{}

			rf := &handlersfakes.FakeResolverFactory{}
			rf.NewCompileCReturnsOnCall(0, compileCR)
			rf.NewCompileCCReturnsOnCall(0, compileCCR)

			h := handlers.NewObject(s, rf, "some-project", "dep-1/dep-1.o")
			ctx := pipeline.NewCtx()
			ctx.Nodes = datum.nodes

			var ext string
			if strings.HasSuffix(datum.name, "CC") {
				ext = ".cc"
			} else {
				ext = ".c"
			}

			source := "dep-1/dep-1" + ext
			sourceN := node.Find(source, ctx.Nodes)
			require.NotNil(t, sourceN)

			name := "dep-1/dep-1.o"
			objectN := node.New(name).Dependency(sourceN)
			if strings.HasSuffix(datum.name, "CC") {
				objectN.Resolver = compileCCR
			} else {
				objectN.Resolver = compileCR
			}
			exNodes := append(datum.nodes.Copy().Cast(), objectN)

			assert.Nil(t, h.Handle(ctx))
			assert.Nil(t, deep.Equal(exNodes, ctx.Nodes))
		})
	}
}
