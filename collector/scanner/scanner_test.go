package scanner_test

import (
	"testing"

	"github.com/ankeesler/btool/collector"
	"github.com/ankeesler/btool/collector/scanner"
	"github.com/ankeesler/btool/collector/scanner/scannerfakes"
	"github.com/ankeesler/btool/collector/testutil"
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

	exNodes := makeExNodes(files)
	store := testutil.FakeStore(exNodes...)
	require.Nil(t, s.Produce(store))

	assert.Equal(t, 1, w.WalkCallCount())
	acRoot, acExts := w.WalkArgsForCall(0)
	assert.Equal(t, "root", acRoot)
	assert.Equal(t, []string{"some-exts"}, acExts)

	assert.Equal(t, 5, store.SetCallCount())
	for i := range files {
		assert.Equal(t, exNodes[i], store.SetArgsForCall(i))
	}
}

func makeExNodes(files []string) []*node.Node {
	nodes := make([]*node.Node, 0)
	for _, file := range files {
		n := node.New(file)
		collector.MustToLabels(n, collector.Labels{Local: true})
		nodes = append(nodes, n)
	}
	return nodes
}
