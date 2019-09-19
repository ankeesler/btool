// Package watcher provides a type that listens to changes in a node.Node graph
// on disk.
package watcher

import (
	"os"

	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
	"gopkg.in/fsnotify.v1"
)

// Watcher is a type that listens to changes in a node.Node graph on disk.
type Watcher struct {
}

// New create a new Watcher.
func New() *Watcher {
	return &Watcher{}
}

// Watch will block until one of the node.Node's in the node.Node graph changes
// on disk. Once the node.Node changes on disk, the call will return. An error is
// returned if the call is unable to listen for changes or if there was a problem
// during listening for changes.
func (w *Watcher) Watch(n *node.Node) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return errors.Wrap(err, "new watcher")
	}
	defer watcher.Close()

	if err := node.Visit(n, func(vn *node.Node) error {
		if _, err := os.Lstat(vn.Name); os.IsNotExist(err) {
			// ok cool! we won't watch this file that doesn't exit
		} else if err != nil {
			return errors.Wrap(err, "lstat")
		} else if err := watcher.Add(vn.Name); err != nil {
			return errors.Wrap(err, "add")
		}

		return nil
	}); err != nil {
		return errors.Wrap(err, "visit")
	}

	select {
	case err := <-watcher.Errors:
		return errors.Wrap(err, "watcher error")
	case event := <-watcher.Events:
		log.Debugf("hooray! caught a watcher event: %s", event)
	}

	return nil
}
