package handlers

import (
	"fmt"
	"sort"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/pipeline"
	"github.com/sirupsen/logrus"
)

type sortTopo struct {
}

// NewSortTopo provides a pipeline.Handler that sorts a node.Node list in
// topological order.
func NewSortTopo() pipeline.Handler {
	return &sortTopo{}
}

func (st *sortTopo) Handle(ctx *pipeline.Ctx) error {
	nodes := ctx.Nodes

	logrus.Debugf("sorting %d nodes", len(nodes))

	sorted := make([]*node.Node, 0, len(nodes))
	sortedSet := make(map[*node.Node]bool)

	for len(sorted) != len(nodes) {
		nodesWithoutDependencies := collectNodesWithoutDependencies(nodes, sortedSet)
		logrus.Debug("nodesWithoutDependencies:", nodesWithoutDependencies)

		if len(nodesWithoutDependencies) == 0 {
			return fmt.Errorf("cycle detected, cannot proceed past %v", sortedSet)
		}

		sorted = append(sorted, nodesWithoutDependencies...)
		for _, node := range nodesWithoutDependencies {
			sortedSet[node] = true
		}
	}

	ctx.Nodes = sorted

	return nil
}

func (st *sortTopo) Name() string { return "topo sort" }

func collectNodesWithoutDependencies(
	nodes []*node.Node,
	sorted map[*node.Node]bool,
) []*node.Node {
	nodesWithoutDependencies := make([]*node.Node, 0)
	for _, node := range nodes {
		if _, ok := sorted[node]; ok {
			continue
		}

		withoutDependencies := true
		for _, dependency := range node.Dependencies {
			if _, ok := sorted[dependency]; !ok {
				withoutDependencies = false
				break
			}
		}

		if withoutDependencies {
			nodesWithoutDependencies = append(nodesWithoutDependencies, node)
		}
	}

	sort.Slice(nodesWithoutDependencies, func(i, j int) bool {
		return nodesWithoutDependencies[i].Name < nodesWithoutDependencies[j].Name
	})

	return nodesWithoutDependencies
}
