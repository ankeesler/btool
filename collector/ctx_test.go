package collector_test

import (
	"testing"

	"github.com/ankeesler/btool/collector"
	"github.com/ankeesler/btool/node"
	"github.com/stretchr/testify/assert"
)

func TestCtxIncludePaths(t *testing.T) {
	// a -> b, c
	// b -> c
	// c
	// d -> a
	cN := node.New("c")
	bN := node.New("b").Dependency(cN)
	aN := node.New("a").Dependency(bN, cN)
	dN := node.New("d").Dependency(aN)

	data := []struct {
		name           string
		includePaths   map[*node.Node]string
		n              *node.Node
		exIncludePaths []string
	}{
		{
			name: "Basic",
			includePaths: map[*node.Node]string{
				aN: "a-path",
				bN: "b-c-path",
				cN: "b-c-path",
				dN: "d-path",
			},
			n: aN,
			exIncludePaths: []string{
				"b-c-path",
				"a-path",
			},
		},
	}
	for _, datum := range data {
		ctx := collector.NewCtx(nil, nil)
		t.Run(datum.name, func(t *testing.T) {
			for n, path := range datum.includePaths {
				ctx.SetIncludePath(n, path)
			}
			acIncludePaths := ctx.IncludePaths(datum.n)
			assert.Equal(t, datum.exIncludePaths, acIncludePaths)
		})
	}
}

func TestCtxIncludePath(t *testing.T) {
	n := node.New("some/path/to/header.h")
	ns := collector.NewNodeStore(nil)
	ns.Add(n)
	ctx := collector.NewCtx(ns, nil)
	assert.Equal(t, "", ctx.IncludePath("header.h"))
	ctx.SetIncludePath(n, "some/path/to")
	assert.Equal(t, "some/path/to", ctx.IncludePath("header.h"))
}
