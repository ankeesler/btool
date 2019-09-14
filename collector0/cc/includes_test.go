package cc_test

import (
	"testing"

	"github.com/ankeesler/btool/collector0/cc"
	"github.com/ankeesler/btool/collector0/cc/ccfakes"
	"github.com/ankeesler/btool/collector0/testutil"
	"github.com/ankeesler/btool/node"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIncludesConsume(t *testing.T) {
	iser := &ccfakes.FakeIncludeser{}
	iser.IncludesReturnsOnCall(0, []string{"a.h", "b/b.h", "gtest/gtest.h"}, nil)

	i := cc.NewIncludes(iser)

	ah := node.New("some/path/to/root/a/a.h")
	bh := node.New("some/path/to/root/b/b.h")
	gtesth := node.New("deps/path/gtest/gtest.h")
	ac := node.New("some/path/to/root/a/a.c")
	s := testutil.FakeStore(ah, bh, gtesth, ac)

	require.Nil(t, i.Consume(s, ac))

	assert.Equal(t, 1, iser.IncludesCallCount())
	assert.Equal(t, "some/path/to/root/a/a.c", iser.IncludesArgsForCall(0))

	assert.Equal(t, 1, s.SetCallCount())
	assert.Equal(t, ac.Dependency(ah, bh, gtesth), s.SetArgsForCall(0))
	assert.Equal(t, "some/path/to/root/,deps/path/,", ac.Labels[cc.LabelIncludePaths])
}
