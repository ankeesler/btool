package cc_test

import (
	"testing"

	"github.com/ankeesler/btool/collector/cc"
	"github.com/ankeesler/btool/collector/cc/ccfakes"
	"github.com/ankeesler/btool/collector/testutil"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/nodefakes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestObject(t *testing.T) {
	// mainc -> bh -> ah
	ah := node.New("a.h")
	ah.Labels[cc.LabelIncludePaths] = "some/include/path,some/other/include/path"
	bh := node.New("b.h")
	bh.Dependency(ah)
	bh.Labels[cc.LabelIncludePaths] = "some/other/other/include/path"
	mainc := node.New("main.c")
	mainc.Labels[cc.LabelIncludePaths] = "some/main/include/path"
	mainc.Dependency(bh)
	maincc := node.New("main.cc")
	maincc.Labels[cc.LabelIncludePaths] = "some/main/include/path"
	maincc.Dependency(bh)

	t.Run("BasicC", func(t *testing.T) {
		r := &nodefakes.FakeResolver{}

		rf := &ccfakes.FakeResolverFactory{}
		rf.NewCompileCReturnsOnCall(0, r)
		o := cc.NewObject(rf)

		s := testutil.FakeStore(ah, bh, mainc)
		require.Nil(t, o.Consume(s, mainc))

		assert.Equal(t, 1, rf.NewCompileCCallCount())
		exIncludePaths := []string{
			"some/include/path",
			"some/other/include/path",
			"some/other/other/include/path",
			"some/main/include/path",
		}
		assert.Equal(t, exIncludePaths, rf.NewCompileCArgsForCall(0))

		assert.Equal(t, 1, s.SetCallCount())
		maino := node.New("main.o")
		maino.Dependency(mainc)
		maino.Resolver = r
		assert.Equal(t, maino, s.SetArgsForCall(0))
	})

	t.Run("BasicCC", func(t *testing.T) {
		r := &nodefakes.FakeResolver{}

		rf := &ccfakes.FakeResolverFactory{}
		rf.NewCompileCCReturnsOnCall(0, r)
		o := cc.NewObject(rf)

		s := testutil.FakeStore(ah, bh, mainc)
		require.Nil(t, o.Consume(s, maincc))

		assert.Equal(t, 1, rf.NewCompileCCCallCount())
		exIncludePaths := []string{
			"some/include/path",
			"some/other/include/path",
			"some/other/other/include/path",
			"some/main/include/path",
		}
		assert.Equal(t, exIncludePaths, rf.NewCompileCCArgsForCall(0))

		assert.Equal(t, 1, s.SetCallCount())
		maino := node.New("main.o")
		maino.Dependency(maincc)
		maino.Resolver = r
		assert.Equal(t, maino, s.SetArgsForCall(0))
	})

	t.Run("NonCCFile", func(t *testing.T) {
		rf := &ccfakes.FakeResolverFactory{}
		o := cc.NewObject(rf)

		s := testutil.FakeStore(ah, bh, mainc)
		require.Nil(t, o.Consume(s, node.New("a.h")))
		assert.Equal(t, 0, rf.NewCompileCCallCount())
		assert.Equal(t, 0, s.SetCallCount())
	})
}
