package cc_test

import (
	"testing"

	"github.com/ankeesler/btool/collector/cc"
	"github.com/ankeesler/btool/collector/collectorfakes"
	"github.com/ankeesler/btool/collector/sorter"
	"github.com/ankeesler/btool/collector/testutil"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/nodefakes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExe(t *testing.T) {
	ah := node.New("a/a.h")
	ac := node.New("a/a.c").Dependency(ah)
	acc := node.New("a/a.cc").Dependency(ah)
	ao := node.New("a/a.o").Dependency(ac)
	bh := node.New("b/b.h").Dependency(ah).Label(cc.LabelLibraries, "c.a,")
	bc := node.New("b/b.c").Dependency(ah, bh)
	bcc := node.New("b/b.cc").Dependency(ah, bh)
	bo := node.New("b/b.o").Dependency(bc)
	mainc := node.New("main.c").Dependency(bh)
	maincc := node.New("main.cc").Dependency(bh)
	maino := node.New("main.o").Dependency(mainc)
	ca := node.New("c.a")

	t.Run("C", func(t *testing.T) {
		r := &nodefakes.FakeResolver{}

		rf := &collectorfakes.FakeResolverFactory{}
		rf.NewLinkCReturnsOnCall(0, r)

		e := cc.NewExe(rf)

		s := testutil.FakeStore(maino, mainc, ao, ac, ah, bo, bc, bh, ca)

		main := node.New("main")
		require.Nil(t, e.Consume(s, main))

		assert.Equal(t, 1, rf.NewLinkCCallCount())

		assert.Equal(t, 1, s.SetCallCount())
		exMain := node.New("main").Dependency(maino, bo, ao, ca)
		exMain.Resolver = r
		sorter.New().Sort(exMain)
		assert.Equal(t, exMain, s.SetArgsForCall(0))
	})

	t.Run("CC", func(t *testing.T) {
		r := &nodefakes.FakeResolver{}

		rf := &collectorfakes.FakeResolverFactory{}
		rf.NewLinkCCReturnsOnCall(0, r)

		e := cc.NewExe(rf)

		s := testutil.FakeStore(maino, maincc, ao, acc, ah, bo, bcc, bh, ca)

		main := node.New("main")
		require.Nil(t, e.Consume(s, main))

		assert.Equal(t, 1, rf.NewLinkCCCallCount())

		assert.Equal(t, 1, s.SetCallCount())
		exMain := node.New("main").Dependency(maino, bo, ao, ca)
		exMain.Resolver = r
		sorter.New().Sort(exMain)
		assert.Equal(t, exMain, s.SetArgsForCall(0))
	})

	t.Run("LoneHeader", func(t *testing.T) {
		r := &nodefakes.FakeResolver{}

		rf := &collectorfakes.FakeResolverFactory{}
		rf.NewLinkCCReturnsOnCall(0, r)

		e := cc.NewExe(rf)

		masterh := node.New("master.h")
		maincc.Dependency(masterh)
		s := testutil.FakeStore(maino, maincc, ao, acc, ah, bo, bcc, bh, masterh, ca)

		main := node.New("main")
		require.Nil(t, e.Consume(s, main))

		assert.Equal(t, 1, rf.NewLinkCCCallCount())

		assert.Equal(t, 1, s.SetCallCount())
		exMain := node.New("main").Dependency(maino, bo, ao, ca)
		exMain.Resolver = r
		sorter.New().Sort(exMain)
		assert.Equal(t, exMain, s.SetArgsForCall(0))
	})

	t.Run("Noop", func(t *testing.T) {
		rf := &collectorfakes.FakeResolverFactory{}
		e := cc.NewExe(rf)
		s := testutil.FakeStore()
		require.Nil(t, e.Consume(s, node.New("main.c")))
		assert.Equal(t, 0, rf.NewLinkCCallCount())
		assert.Equal(t, 0, s.SetCallCount())
	})
}
