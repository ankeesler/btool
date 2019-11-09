package cc_test

import (
	"testing"

	"github.com/ankeesler/btool/app/collector/cc"
	"github.com/ankeesler/btool/app/collector/collectorfakes"
	"github.com/ankeesler/btool/app/collector/testutil"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/nodefakes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestObj(t *testing.T) {
	// mainc -> bh -> ah
	ah := node.New("a.h")
	require.Nil(t, cc.AppendIncludePaths(
		ah,
		"some/include/path",
		"some/other/include/path",
	))
	bh := node.New("b.h")
	bh.Dependency(ah)
	require.Nil(t, cc.AppendIncludePaths(
		bh,
		"some/other/other/include/path",
	))
	mainc := node.New("main.c")
	mainc.Dependency(bh)
	require.Nil(t, cc.AppendIncludePaths(mainc,
		"some/main/include/path",
	))
	maincc := node.New("main.cc")
	maincc.Dependency(bh)
	require.Nil(t, cc.AppendIncludePaths(maincc,
		"some/main/include/path",
	))

	t.Run("BasicC", func(t *testing.T) {
		r := &nodefakes.FakeResolver{}

		rf := &collectorfakes.FakeResolverFactory{}
		rf.NewCompileCReturnsOnCall(0, r)
		o := cc.NewObj(rf)

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

		rf := &collectorfakes.FakeResolverFactory{}
		rf.NewCompileCCReturnsOnCall(0, r)
		o := cc.NewObj(rf)

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
		rf := &collectorfakes.FakeResolverFactory{}
		o := cc.NewObj(rf)

		s := testutil.FakeStore(ah, bh, mainc)
		require.Nil(t, o.Consume(s, node.New("a.h")))
		assert.Equal(t, 0, rf.NewCompileCCallCount())
		assert.Equal(t, 0, s.SetCallCount())
	})

	t.Run("FileWithDotC", func(t *testing.T) {
		r := &nodefakes.FakeResolver{}

		rf := &collectorfakes.FakeResolverFactory{}
		rf.NewCompileCReturnsOnCall(0, r)
		o := cc.NewObj(rf)

		mainc := node.New("github.com/ankeesler/btool/example/BasicC/main.c")
		s := testutil.FakeStore(mainc)
		require.Nil(t, o.Consume(s, mainc))

		assert.Equal(t, 1, s.SetCallCount())
		maino := node.New("github.com/ankeesler/btool/example/BasicC/main.o")
		maino.Dependency(mainc)
		maino.Resolver = r
		assert.Equal(t, maino, s.SetArgsForCall(0))
	})
}
