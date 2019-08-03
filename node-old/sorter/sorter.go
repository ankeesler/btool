// Package sorter provides a node.Handler that performs a topological sort on nodes.
package sorter

import (
	"errors"
	"sort"

	"github.com/ankeesler/btool/node"
	"github.com/sirupsen/logrus"
)

type Sorter struct {
}

func New() *Sorter {
	return &Sorter{}
}

func (s *Sorter) Handle(nodes []*node.Node) ([]*node.Node, error) {
	logrus.Debugf("sorting %d nodes", len(nodes))

	sorted := make([]*node.Node, 0, len(nodes))
	sortedSet := make(map[*node.Node]bool)

	for len(sorted) != len(nodes) {
		nodesWithoutDependencies := collectNodesWithoutDependencies(nodes, sortedSet)
		logrus.Debug("nodesWithoutDependencies:", nodesWithoutDependencies)

		if len(nodesWithoutDependencies) == 0 {
			return nil, errors.New("cycle detected")
		}

		sorted = append(sorted, nodesWithoutDependencies...)
		for _, node := range nodesWithoutDependencies {
			sortedSet[node] = true
		}
	}

	return sorted, nil
}

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
