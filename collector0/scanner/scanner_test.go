package scanner_test

import (
	"testing"

	collector "github.com/ankeesler/btool/collector0"
	"github.com/ankeesler/btool/collector0/scanner"
	"github.com/ankeesler/btool/collector0/scanner/scannerfakes"
	"github.com/ankeesler/btool/node"
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

	s := scanner.New(w, "root", []string{"some-exts"})

	acStore := collector.NewStore()
	require.Nil(t, s.Produce(acStore))

	exStore := collector.NewStore()
	exNodes := makeExNodes(files)
	for _, exN := range exNodes {
		exStore.Add(exN)
	}
	assert.Equal(t, exStore, acStore)

	assert.Equal(t, 1, w.WalkCallCount())
	acRoot, acExts := w.WalkArgsForCall(0)
	assert.Equal(t, "root", acRoot)
	assert.Equal(t, []string{"some-exts"}, acExts)
}

func makeExNodes(files []string) []*node.Node {
	nodes := make([]*node.Node, 0)
	for _, file := range files {
		nodes = append(nodes, node.New(file))
	}
	return nodes
}
