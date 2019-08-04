package sorter_test

import (
	"reflect"
	"testing"

	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/sorter"
	"github.com/sirupsen/logrus"
)

func TestAlphaHandle(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

	s := sorter.NewAlpha()
	cfg := node.Config{}

	nodeA := node.New("a")
	nodeB := node.New("b")
	nodeC := node.New("c")

	node0 := node.New("0")
	node1 := node.New("1")
	node2 := node.New("2")
	node3 := node.New("3")

	nodeACpy := *nodeA
	nodeBCpy := *nodeB
	nodeCCpy := *nodeC

	nodes := []*node.Node{
		&nodeBCpy,
		nodeACpy.Dependency(node1).Dependency(node0),
		nodeCCpy.Dependency(node3).Dependency(node2).Dependency(node1),
	}

	ac, err := s.Handle(&cfg, nodes)
	if err != nil {
		t.Error(err)
	}

	ex := []*node.Node{
		nodeA.Dependency(node0).Dependency(node1),
		nodeB,
		nodeC.Dependency(node1).Dependency(node2).Dependency(node3),
	}
	if !reflect.DeepEqual(ex, ac) {
		t.Error(ex, "!=", ac)
	}
}
