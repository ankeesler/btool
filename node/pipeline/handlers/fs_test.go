package handlers_test

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node/pipeline"
	"github.com/ankeesler/btool/node/pipeline/handlers"
	"github.com/ankeesler/btool/node/pipeline/handlers/handlersfakes"
	"github.com/ankeesler/btool/node/testutil"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFS(t *testing.T) {
	data := []struct {
		name    string
		exNodes testutil.Nodes
	}{
		{
			name:    "BasicC",
			exNodes: testutil.BasicNodesC.Copy(),
		},
		{
			name:    "BasicCC",
			exNodes: testutil.BasicNodesCC.Copy(),
		},
	}
	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			root := "/"
			c := &handlersfakes.FakeCollector{}
			c.CollectReturnsOnCall(0, paths(root, datum.exNodes), nil)
			i := &handlersfakes.FakeIncludeser{}
			setupIncludeserFake(root, i, datum.exNodes)
			h := handlers.NewFS(c, i, root)
			ctx := pipeline.NewCtx()
			err := h.Handle(ctx)
			require.Nil(t, err)

			exNodes := datum.exNodes
			for _, exN := range exNodes {
				exN.Name = filepath.Join("/", exN.Name)
			}

			assert.Nil(t, deep.Equal(datum.exNodes.Cast(), ctx.Nodes))
		})
	}
}

func paths(root string, nodes testutil.Nodes) []string {
	p := make([]string, 0)
	for _, n := range nodes {
		p = append(p, filepath.Join(root, n.Name))
	}
	return p
}

func setupIncludeserFake(
	root string,
	i *handlersfakes.FakeIncludeser,
	nodes testutil.Nodes,
) {
	includes := make(map[string][]string)
	for _, n := range nodes {
		name := filepath.Join(root, n.Name)
		if _, ok := includes[name]; !ok {
			includes[name] = make([]string, 0)
		}
		for _, d := range n.Dependencies {
			dName := filepath.Join(root, d.Name)
			if strings.HasSuffix(d.Name, ".h") {
				includes[name] = append(includes[name], dName)
				log.Debugf("includes %s -> %s", name, dName)
			}
		}
	}

	i.IncludesStub = func(path string) ([]string, error) {
		tuna, ok := includes[path]
		if !ok {
			return nil, errors.New("unknown path: " + path)
		} else {
			return tuna, nil
		}
	}
}
