package app_test

import (
	"testing"

	"github.com/ankeesler/btool/app"
	"github.com/ankeesler/btool/app/appfakes"
	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
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
		watch          bool
		cleanCallCount int
		listCallCount  int
		buildCallCount int
		runCallCount   int
		watchCallCount int
		errMsg         string
	}{
		{
			name:           "Build",
			n:              node.New(""),
			clean:          false,
			list:           false,
			run:            false,
			watch:          false,
			cleanCallCount: 0,
			listCallCount:  0,
			buildCallCount: 1,
			runCallCount:   0,
			watchCallCount: 0,
			errMsg:         "",
		},
		{
			name:           "Clean",
			n:              node.New(""),
			clean:          true,
			list:           false,
			run:            false,
			watch:          false,
			cleanCallCount: 1,
			listCallCount:  0,
			buildCallCount: 0,
			runCallCount:   0,
			watchCallCount: 0,
			errMsg:         "",
		},
		{
			name:           "List",
			n:              node.New(""),
			clean:          false,
			list:           true,
			run:            false,
			watch:          false,
			cleanCallCount: 0,
			listCallCount:  1,
			buildCallCount: 0,
			runCallCount:   0,
			watchCallCount: 0,
			errMsg:         "",
		},
		{
			name:           "Run",
			n:              node.New(""),
			clean:          false,
			list:           false,
			run:            true,
			watch:          false,
			cleanCallCount: 0,
			listCallCount:  0,
			buildCallCount: 1,
			runCallCount:   1,
			watchCallCount: 0,
			errMsg:         "",
		},
		{
			name:           "Watch",
			n:              node.New(""),
			clean:          false,
			list:           false,
			run:            false,
			watch:          true,
			cleanCallCount: 0,
			listCallCount:  0,
			buildCallCount: 3,
			runCallCount:   0,
			watchCallCount: 3,
			errMsg:         "watch: some watch error",
		},
		{
			name:           "WatchRun",
			n:              node.New(""),
			clean:          false,
			list:           false,
			run:            true,
			watch:          true,
			cleanCallCount: 0,
			listCallCount:  0,
			buildCallCount: 3,
			runCallCount:   2,
			watchCallCount: 3,
			errMsg:         "watch: some watch error",
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
			builder.BuildReturnsOnCall(1, errors.New("some build error"))
			runner := &appfakes.FakeRunner{}
			runner.RunReturnsOnCall(1, errors.New("some run error"))
			watcher := &appfakes.FakeWatcher{}
			watcher.WatchReturnsOnCall(2, errors.New("some watch error"))

			b := app.New(cc, cleaner, lister, builder, runner, watcher)
			err := b.Run(datum.n, datum.clean, datum.list, datum.run, datum.watch)
			if datum.errMsg == "" {
				require.Nil(t, err)
			} else {
				require.EqualError(t, err, datum.errMsg)
			}

			assert.Equal(t, 1, cc.CreateCallCount())
			assert.Equal(t, 1, c.CollectCallCount())

			require.Equal(t, datum.cleanCallCount, cleaner.CleanCallCount())
			if datum.cleanCallCount > 0 {
				assert.Equal(t, datum.n, cleaner.CleanArgsForCall(0))
			}

			require.Equal(t, datum.buildCallCount, builder.BuildCallCount())
			if datum.buildCallCount > 0 {
				assert.Equal(t, datum.n, builder.BuildArgsForCall(0))
			}

			require.Equal(t, datum.runCallCount, runner.RunCallCount())
			if datum.runCallCount > 0 {
				assert.Equal(t, datum.n, runner.RunArgsForCall(0))
			}

			require.Equal(t, datum.watchCallCount, watcher.WatchCallCount())
			if datum.watchCallCount > 0 {
				assert.Equal(t, datum.n, watcher.WatchArgsForCall(0))
			}
		})
	}
}
