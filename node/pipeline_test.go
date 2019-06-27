package node_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/nodefakes"
	"github.com/ankeesler/btool/node/testutil"
	"github.com/sirupsen/logrus"
)

func TestPipeline(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

	goodHandlerA := &nodefakes.FakeHandler{}
	goodHandlerA.HandleReturnsOnCall(0, testutil.BasicNodes, nil)
	goodHandlerB := &nodefakes.FakeHandler{}
	goodHandlerB.HandleReturnsOnCall(0, testutil.BasicNodes, nil)
	badHandler := &nodefakes.FakeHandler{}
	badHandler.HandleReturnsOnCall(0, nil, errors.New("some error"))

	// Happy.
	if err := node.Pipeline(goodHandlerA, goodHandlerB); err != nil {
		t.Error(err)
	}

	// Sad.
	if err := node.Pipeline(goodHandlerA, badHandler, goodHandlerB); err == nil {
		t.Error("expected failure")
	} else if err.Error() != "some error" {
		t.Error("expected 'some error'")
	}

	data := []struct {
		ex int
		f  func() int
	}{
		{
			ex: 2,
			f:  goodHandlerA.HandleCallCount,
		},
		{
			ex: 1,
			f:  goodHandlerB.HandleCallCount,
		},
		{
			ex: 1,
			f:  badHandler.HandleCallCount,
		},
	}
	for _, datum := range data {
		if ac := datum.f(); datum.ex != ac {
			t.Errorf("%s: %d != %d", reflect.TypeOf(datum.f).Name(), datum.ex, ac)
		}
	}
}
