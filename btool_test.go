package btool_test

import (
	"testing"

	"github.com/ankeesler/btool"
	"github.com/ankeesler/btool/btoolfakes"
	"github.com/ankeesler/btool/node"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBtoolRun(t *testing.T) {
	data := []struct {
		name             string
		n                *node.Node
		clean            bool
		collectCallCount int
		cleanCallCount   int
		buildCallCount   int
	}{
		{
			name:             "Build",
			n:                node.New(""),
			clean:            false,
			collectCallCount: 1,
			cleanCallCount:   0,
			buildCallCount:   1,
		},
		{
			name:             "Clean",
			n:                node.New(""),
			clean:            true,
			collectCallCount: 1,
			cleanCallCount:   1,
			buildCallCount:   0,
		},
	}
	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			collector := &btoolfakes.FakeCollector{}
			collector.CollectReturnsOnCall(0, nil)

			cleaner := &btoolfakes.FakeCleaner{}

			builder := &btoolfakes.FakeBuilder{}

			b := btool.New(collector, cleaner, builder)
			require.Nil(t, b.Run(datum.n, datum.clean, false))

			assert.Equal(t, datum.collectCallCount, collector.CollectCallCount())
			if datum.collectCallCount > 0 {
				assert.Equal(t, datum.n, collector.CollectArgsForCall(0))
			}

			assert.Equal(t, datum.cleanCallCount, cleaner.CleanCallCount())
			if datum.cleanCallCount > 0 {
				assert.Equal(t, datum.n, cleaner.CleanArgsForCall(0))
			}

			assert.Equal(t, datum.buildCallCount, builder.BuildCallCount())
			if datum.buildCallCount > 0 {
				assert.Equal(t, datum.n, builder.BuildArgsForCall(0))
			}
		})
	}
}
