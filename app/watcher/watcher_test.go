package watcher_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ankeesler/btool/app/watcher"
	"github.com/ankeesler/btool/node"
	"github.com/stretchr/testify/require"
)

func TestWatcherWatch(t *testing.T) {
	dir, err := ioutil.TempDir("", "btool_watcher_test")
	require.Nil(t, err)
	defer func() {
		require.Nil(t, os.RemoveAll(dir))
	}()

	w := watcher.New()

	// a -> b, c
	// b -> c
	// c -> d
	// d
	d := node.New(filepath.Join(dir, "d"))
	c := node.New(filepath.Join(dir, "c")).Dependency(d)
	b := node.New(filepath.Join(dir, "b")).Dependency(c)
	a := node.New(filepath.Join(dir, "a")).Dependency(b, c)
	node.Visit(a, func(vn *node.Node) error {
		if !strings.HasSuffix(vn.Name, "d") {
			require.Nil(t, ioutil.WriteFile(vn.Name, []byte(vn.Name), 0644))
		}
		return nil
	})

	done := make(chan struct{})
	go func() {
		for {
			require.Nil(t, w.Watch(a))
			done <- struct{}{}
		}
	}()

	assertEmpty(t, done)

	require.Nil(t, ioutil.WriteFile(a.Name, []byte("a"), 0644))

	assertNotEmpty(t, done)
	assertEmpty(t, done)

	require.Nil(t, ioutil.WriteFile(d.Name, []byte("d"), 0644))
	require.Nil(t, ioutil.WriteFile(c.Name, []byte("c"), 0644))
	require.Nil(t, ioutil.WriteFile(b.Name, []byte("b"), 0644))

	assertNotEmpty(t, done)
	assertEmpty(t, done)

	require.Nil(t, ioutil.WriteFile(c.Name, []byte("c"), 0644))
	require.Nil(t, ioutil.WriteFile(b.Name, []byte("b"), 0644))

	assertNotEmpty(t, done)
	assertEmpty(t, done)
}

func assertEmpty(t *testing.T, doneC <-chan struct{}) {
	timer := time.NewTimer(time.Millisecond * 500)
	select {
	case <-timer.C:
		// yay!
	case <-doneC:
		if !timer.Stop() {
			<-timer.C
		}
		require.Fail(t, "did not expect call to return")
	}
}

func assertNotEmpty(t *testing.T, doneC <-chan struct{}) {
	timer := time.NewTimer(time.Millisecond * 500)
	select {
	case <-timer.C:
		require.Fail(t, "expected ")
	case <-doneC:
		// yay!
		if !timer.Stop() {
			<-timer.C
		}
	}
}
