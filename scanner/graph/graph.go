// Package graph provides a simple directed graph data structure.
package graph

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/pkg/errors"
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

// Topological sort.
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

// Walk the graph in a "breadth" first way.
func (g *Graph) Walk(root string, handler func(node *Node) error) error {
	if g.nodeNames[root] == nil {
		return errors.New("unknown node name " + root)
	}
	return g.walk(root, handler, make(map[string]bool))
}

func (g *Graph) String() string {
	buf := bytes.NewBuffer([]byte{})

	for _, nodeName := range sortKeys0(g.nodes) {
		dependenciesNames := g.nodes[nodeName]
		buf.WriteString(fmt.Sprintf("%s:\n", nodeName))

		for _, dependencyName := range sortKeys1(dependenciesNames) {
			buf.WriteString(fmt.Sprintf("> %s\n", dependencyName))
		}
	}

	return buf.String()
}

func (g *Graph) Edges(node *Node) []*Node {
	edgesMap, ok := g.nodes[node.Name]
	if !ok {
		return nil
	}

	edges := make([]*Node, 0, len(edgesMap))
	for edgeName := range edgesMap {
		edges = append(edges, g.nodeNames[edgeName])
	}
	return edges
}

func Equal(left, right *Graph) error {
	if err := superset(left, right); err != nil {
		return errors.Wrap(err, "left -> right")
	} else if err := superset(right, left); err != nil {
		return errors.Wrap(err, "right -> left")
	} else {
		return nil
	}
}

func superset(left, right *Graph) error {
	for _, nodeName := range sortKeys0(left.nodes) {
		dependenciesNames := left.nodes[nodeName]
		otherDependenciesNames, ok := right.nodes[nodeName]
		if !ok {
			return fmt.Errorf("node %s does not exist", nodeName)
		}

		for _, dependencyName := range sortKeys1(dependenciesNames) {
			_, ok := otherDependenciesNames[dependencyName]
			if !ok {
				return fmt.Errorf(
					"node %s is missing dependency %s",
					nodeName,
					dependencyName,
				)
			}
		}
	}
	return nil
}

func sortKeys0(m map[string]map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for key, _ := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func sortKeys1(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for key, _ := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (g *Graph) walk(
	root string,
	handler func(node *Node) error,
	visited map[string]bool,
) error {
	if visited[root] {
		return nil
	}

	if err := handler(g.nodeNames[root]); err != nil {
		return errors.Wrap(err, fmt.Sprintf("handle %s", root))
	}

	visited[root] = true

	for dependency, _ := range g.nodes[root] {
		if err := g.walk(dependency, handler, visited); err != nil {
			return err
		}
	}

	return nil
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
	newNodes := make([]*Node, 0, 2)
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
			newNodes = append(newNodes, g.nodeNames[nodeName])
		}
	}

	sort.Slice(newNodes, func(i, j int) bool {
		return newNodes[i].Name < newNodes[j].Name
	})
	return append(nodesWithoutDependencies, newNodes...)
}

func nodeNameInSlice(nodeName string, slice []*Node) bool {
	for _, otherNode := range slice {
		if nodeName == otherNode.Name {
			return true
		}
	}
	return false
}
