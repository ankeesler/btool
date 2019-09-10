package scanner_test

import (
	"testing"

	"github.com/ankeesler/btool/collector/prodcon"
	"github.com/ankeesler/btool/collector/scanner"
	"github.com/ankeesler/btool/node"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScannerProduce(t *testing.T) {
	w := &scannerfakes.FakeWalker{}
	files := []string{
		"main.c",
		"a/a.c",
		"a/a.h",
		"b/b.c",
		"b/b.h",
	}
	w.WalkReturnsOnCall(0, files, nil)

	s := scanner.New(w)

	acStore := prodcon.NewStore()
	require.Nil(t, s.Produce(acStore))

	exNodes := makeExNodes(files)
	for _, exN := range exNodes {
		exStore.Add(exN)
	}
	assert.Nil(t, deep.Equal(exStore, acStore))
}

func makeExNodes(files []string) {
	nodes := make([]*node.Node, 0)
	for _, file := range files {
		nodes = append(nodes, node.New(file))
	}
	return nodes
}
