package resolver

import (
	"github.com/ankeesler/btool/node"
	"github.com/sirupsen/logrus"
)

type NoopResolver struct {
}

func NewNoop() *NoopResolver {
	return &NoopResolver{}
}

func (nr *NoopResolver) Handle(cfg *node.Config, nodes []*node.Node) ([]*node.Node, error) {
	for _, n := range nodes {
		n.Resolver = nr
	}
	return nodes, nil
}

func (nr *NoopResolver) Resolve(n *node.Node) error {
	logrus.Infof("resolve %s", n)
	return nil
}
