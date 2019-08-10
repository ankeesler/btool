package deps

import (
	"fmt"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/pipeline"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Deps struct {
	Deps []Dep `yaml: "deps"`
}

func (d *Deps) Nodes(ctx *pipeline.Ctx) error {
	rf := newResolverFactory(ctx)
	for _, dep := range d.Deps {
		var err error
		ctx.Nodes, err = dep.nodes(ctx.Nodes, rf)
		if err != nil {
			return err
		}
	}
	return nil
}

type Dep struct {
	Name   string `yaml: "name"`
	URL    string `yaml: "url"`
	SHA256 string `yaml: "sha256"`
	Nodes  []Node `yaml: "node"`
}

func (d *Dep) nodes(nodes []*node.Node, rf *resolverFactory) ([]*node.Node, error) {
	for _, depNode := range d.Nodes {
		depN := node.New(depNode.Name)
		nodes = append(nodes, depN)

		if depNode.Resolver != "" {
			logrus.Debugf("resolver %s for node %s", depNode.Resolver, depNode.Name)
			r := rf.make(depNode.Resolver)
			if r == nil {
				return nodes, fmt.Errorf("unknown resolver: %s", depNode.Resolver)
			}

			depN.Resolver = r
		}
	}

	for _, depNode := range d.Nodes {
		depN := node.Find(depNode.Name, nodes)
		if depN == nil {
			return nodes, errors.New("couldn't find node " + depNode.Name)
		}

		for _, depDependency := range depNode.Dependencies {
			depDependencyN := node.Find(depDependency, nodes)
			if depDependencyN == nil {
				return nodes, fmt.Errorf(
					"couldn't find dependency %s for node %s",
					depDependency,
					depNode.Name,
				)
			}

			depN.Dependency(depDependencyN)
		}
	}

	return nodes, nil
}

type Node struct {
	Name         string   `yaml: "name"`
	Dependencies []string `yaml: "dependencies"`
	Resolver     string   `yaml: "resolver"`
}
