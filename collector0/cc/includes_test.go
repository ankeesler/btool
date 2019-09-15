package cc_test

import (
	"testing"

	collector "github.com/ankeesler/btool/collector0"
	"github.com/ankeesler/btool/collector0/cc"
	"github.com/ankeesler/btool/collector0/cc/ccfakes"
	"github.com/ankeesler/btool/collector0/testutil"
	"github.com/ankeesler/btool/node"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIncludesConsume(t *testing.T) {
	t.Run("C", func(t *testing.T) {
		iser := &ccfakes.FakeIncludeser{}
		iser.IncludesReturnsOnCall(0, []string{"a.h", "b/b.h", "gtest/gtest.h"}, nil)

		i := cc.NewIncludes(iser)

		ah := node.New("some/path/to/root/a/a.h")
		bh := node.New("some/path/to/root/b/b.h")
		gtesth := node.New("deps/path/gtest/gtest.h")
		ac := node.New("some/path/to/root/a/a.c").Label(collector.LabelLocal, "true")
		s := testutil.FakeStore(ah, bh, gtesth, ac)

		require.Nil(t, i.Consume(s, ac))

		assert.Equal(t, 1, iser.IncludesCallCount())
		assert.Equal(t, "some/path/to/root/a/a.c", iser.IncludesArgsForCall(0))

		assert.Equal(t, 1, s.SetCallCount())
		assert.Equal(t, ac.Dependency(ah, bh, gtesth), s.SetArgsForCall(0))
		assert.Equal(t, "some/path/to/root/,deps/path/,", ac.Labels[cc.LabelIncludePaths])
	})

	t.Run("NotLocal", func(t *testing.T) {
		iser := &ccfakes.FakeIncludeser{}
		i := cc.NewIncludes(iser)

		s := testutil.FakeStore()
		ac := node.New("a.c")
		require.Nil(t, i.Consume(s, ac))
		assert.Equal(t, 0, iser.IncludesCallCount())
		assert.Equal(t, 0, s.SetCallCount())
	})

	t.Run("BadFileExt", func(t *testing.T) {
		iser := &ccfakes.FakeIncludeser{}
		i := cc.NewIncludes(iser)

		s := testutil.FakeStore()
		azip := node.New("a.zip").Label(collector.LabelLocal, "true")
		require.Nil(t, i.Consume(s, azip))
		assert.Equal(t, 0, iser.IncludesCallCount())
		assert.Equal(t, 0, s.SetCallCount())
	})
}
