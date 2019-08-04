package handlers

import (
	"fmt"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/pipeline"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type resolve struct{}

// NewResolve returns a pipeline.Handler that runs all of the node.Resolver graph
// for a particular target.
func NewResolve() pipeline.Handler {
	return &resolve{}
}

func (r *resolve) Handle(ctx *pipeline.Ctx) {
	target := ctx.KV[pipeline.CtxTarget]
	nodes := ctx.Nodes

	n := node.Find(target, nodes)
	if n == nil {
		ctx.Err = fmt.Errorf("unknown target %s", target)
		return
	}

	if err := build(n); err != nil {
		ctx.Err = errors.Wrap(err, "build")
		return
	}
}

func (r *resolve) Name() string { return "resolve" }

func build(n *node.Node) error {
	logrus.Debugf("building %s", n.Name)

	for _, d := range n.Dependencies {
		logrus.Debugf("building dependency %s", d.Name)
		if err := build(d); err != nil {
			return errors.Wrap(err, "build "+d.Name)
		}
	}

	logrus.Debugf("resolving %s", n.Name)
	if n.Resolver != nil {
		if err := n.Resolver.Resolve(n); err != nil {
			return errors.Wrap(err, "resolve "+n.Name)
		}
	}

	return nil
}
