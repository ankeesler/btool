// Package graph provides a simple directed graph data structure.
package graph

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
)

type Graph struct {
	nodeNames map[string]*Node
	nodes     map[string]map[string]bool
}

func New() *Graph {
	return &Graph{
		nodeNames: make(map[string]*Node),
		nodes:     make(map[string]map[string]bool),
	}
}

func (g *Graph) Add(node, dependency *Node) *Graph {
	logrus.Debugf("graph: add %s -> %s", node, dependency)

	g.nodeNames[node.Name] = node

	if _, ok := g.nodes[node.Name]; !ok {
		g.nodes[node.Name] = make(map[string]bool)
	}

	if dependency != nil {
		g.nodeNames[dependency.Name] = dependency
		g.Add(dependency, nil)

		g.nodes[node.Name][dependency.Name] = true
	}

	return g
}

func (g *Graph) Sort() ([]*Node, error) {
	logrus.Debugf("graph: sorting %d nodes", len(g.nodes))

	g.resetForSort()

	toBeSorted := make([]*Node, 0, len(g.nodes))
	sorted := make([]*Node, 0, len(g.nodes))

	for i := 0; i < len(g.nodes); i++ {
		toBeSorted = g.collectNodesWithoutDependencies(toBeSorted)
		logrus.Debugf("graph: to be sorted: %s", toBeSorted)

		if i >= len(toBeSorted) {
			return nil, errors.New("cycle detected")
		}

		node := toBeSorted[i]

		sorted = append(sorted, node)
		g.updateForSortedNode(node)
	}

	return sorted, nil
}

func (g *Graph) String() string {
	buf := bytes.NewBuffer([]byte{})

	for nodeName, dependenciesNames := range g.nodes {
		buf.WriteString(fmt.Sprintf("%s:\n", nodeName))
		for dependencyName, _ := range dependenciesNames {
			buf.WriteString(fmt.Sprintf("> %s\n", dependencyName))
		}
	}

	return buf.String()
}

func (g *Graph) updateForSortedNode(sortedNode *Node) {
	for _, dependenciesNames := range g.nodes {
		if _, ok := dependenciesNames[sortedNode.Name]; ok {
			dependenciesNames[sortedNode.Name] = false
		}
	}
}

func (g *Graph) resetForSort() {
	for _, dependenciesNames := range g.nodes {
		for dependencyName := range dependenciesNames {
			dependenciesNames[dependencyName] = true
		}
	}
}

func (g *Graph) collectNodesWithoutDependencies(nodesWithoutDependencies []*Node) []*Node {
	for nodeName, dependencies := range g.nodes {
		if nodeNameInSlice(nodeName, nodesWithoutDependencies) {
			continue
		}

		withoutDependencies := true

		for _, inUse := range dependencies {
			if inUse {
				withoutDependencies = false
			}
		}

		if withoutDependencies {
			nodesWithoutDependencies = append(nodesWithoutDependencies, g.nodeNames[nodeName])
		}
	}

	return nodesWithoutDependencies
}

func nodeNameInSlice(nodeName string, slice []*Node) bool {
	for _, otherNode := range slice {
		if nodeName == otherNode.Name {
			return true
		}
	}
	return false
}
