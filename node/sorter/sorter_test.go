package sorter_test

import (
	"reflect"
	"testing"

	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/sorter"
	"github.com/ankeesler/btool/node/testutil"
	"github.com/sirupsen/logrus"
)

func TestHandle(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

	s := sorter.New()

	// Happy.
	ac, err := s.Handle([]*node.Node{
		&testutil.Mainc,
		&testutil.Dep1h,
		&testutil.Dep0h,
	})
	if err != nil {
		t.Error(err)
	}

	ex := []*node.Node{
		&testutil.Dep0h,
		&testutil.Dep1h,
		&testutil.Mainc,
	}
	if !reflect.DeepEqual(ex, ac) {
		t.Error(ex, "!=", ac)
	}

	// Sad.
	testutil.Dep0h.Dependencies = []*node.Node{&testutil.Mainc}
	ac, err = s.Handle(testutil.BasicNodes)
	if err == nil {
		t.Error("expected failure")
	}
}
