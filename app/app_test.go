package app_test

import (
	"testing"

	"github.com/ankeesler/btool/app"
	"github.com/ankeesler/btool/app/appfakes"
	"github.com/ankeesler/btool/node"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBtoolRun(t *testing.T) {
	data := []struct {
		name           string
		n              *node.Node
		clean          bool
		list           bool
		run            bool
		cleanCallCount int
		listCallCount  int
		buildCallCount int
		runCallCount   int
	}{
		{
			name:           "Build",
			n:              node.New(""),
			clean:          false,
			list:           false,
			run:            false,
			cleanCallCount: 0,
			listCallCount:  0,
			buildCallCount: 1,
			runCallCount:   0,
		},
		{
			name:           "Clean",
			n:              node.New(""),
			clean:          true,
			list:           false,
			run:            false,
			cleanCallCount: 1,
			listCallCount:  0,
			buildCallCount: 0,
			runCallCount:   0,
		},
		{
			name:           "List",
			n:              node.New(""),
			clean:          false,
			list:           true,
			run:            false,
			cleanCallCount: 0,
			listCallCount:  1,
			buildCallCount: 0,
			runCallCount:   0,
		},
		{
			name:           "Run",
			n:              node.New(""),
			clean:          false,
			list:           false,
			run:            true,
			cleanCallCount: 0,
			listCallCount:  0,
			buildCallCount: 1,
			runCallCount:   1,
		},
	}
	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			c := &appfakes.FakeCollector{}
			cc := &appfakes.FakeCollectorCreator{}
			cc.CreateReturnsOnCall(0, c, nil)

			cleaner := &appfakes.FakeCleaner{}
			lister := &appfakes.FakeLister{}
			builder := &appfakes.FakeBuilder{}
			runner := &appfakes.FakeRunner{}

			b := app.New(cc, cleaner, lister, builder, runner)
			require.Nil(t, b.Run(datum.n, datum.clean, datum.list, datum.run))

			assert.Equal(t, 1, cc.CreateCallCount())
			assert.Equal(t, 1, c.CollectCallCount())

			assert.Equal(t, datum.cleanCallCount, cleaner.CleanCallCount())
			if datum.cleanCallCount > 0 {
				assert.Equal(t, datum.n, cleaner.CleanArgsForCall(0))
			}

			assert.Equal(t, datum.buildCallCount, builder.BuildCallCount())
			if datum.buildCallCount > 0 {
				assert.Equal(t, datum.n, builder.BuildArgsForCall(0))
			}

			assert.Equal(t, datum.runCallCount, runner.RunCallCount())
			if datum.runCallCount > 0 {
				assert.Equal(t, datum.n, runner.RunArgsForCall(0))
			}
		})
	}
}
