package cc_test

import (
	"testing"

	"github.com/ankeesler/btool/collector"
	"github.com/ankeesler/btool/collector/cc"
	"github.com/ankeesler/btool/collector/cc/ccfakes"
	"github.com/ankeesler/btool/collector/testutil"
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
		ac := node.New("some/path/to/root/a/a.c")
		collector.MustToLabels(ac, collector.Labels{Local: true})
		s := testutil.FakeStore(ah, bh, gtesth, ac)

		require.Nil(t, i.Consume(s, ac))

		labels := cc.Labels{}
		collector.MustFromLabels(ac, &labels)

		assert.Equal(t, 1, iser.IncludesCallCount())
		assert.Equal(t, "some/path/to/root/a/a.c", iser.IncludesArgsForCall(0))

		assert.Equal(t, 1, s.SetCallCount())
		assert.Equal(t, ac.Dependency(ah, bh, gtesth), s.SetArgsForCall(0))
		assert.Equal(
			t,
			[]string{"some/path/to/root/", "deps/path/"},
			labels.IncludePaths,
		)
	})

	t.Run("EmptyIncludePath", func(t *testing.T) {
		iser := &ccfakes.FakeIncludeser{}
		iser.IncludesReturnsOnCall(0, []string{"a/a.h"}, nil)
		i := cc.NewIncludes(iser)

		ah := node.New("a/a.h")
		ac := node.New("a/a.c")
		collector.MustToLabels(ac, collector.Labels{Local: true})
		s := testutil.FakeStore(ah, ac)

		require.Nil(t, i.Consume(s, ac))

		labels := cc.Labels{}
		collector.MustFromLabels(ac, &labels)

		assert.Equal(t, 1, iser.IncludesCallCount())
		assert.Equal(t, "a/a.c", iser.IncludesArgsForCall(0))

		assert.Equal(t, 1, s.SetCallCount())
		assert.Equal(t, ac.Dependency(ah), s.SetArgsForCall(0))
		assert.Equal(t, []string{"."}, labels.IncludePaths)
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
		azip := node.New("a.zip")
		collector.MustToLabels(azip, collector.Labels{Local: true})
		require.Nil(t, i.Consume(s, azip))
		assert.Equal(t, 0, iser.IncludesCallCount())
		assert.Equal(t, 0, s.SetCallCount())
	})
}
