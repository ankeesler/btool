// Package sorter provides a stable way of sorting a node.Node graph.
package sorter

import "sort"

// sortaAlpha sorts the provided node.Node's and dependencies in alphanumeric
// order.
func sortAlpha(nodes []*Node) {
	sohrt(nodes)
	for _, n := range nodes {
		sohrt(n.Dependencies)
	}
}

func sohrt(nodes []*Node) {
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Name < nodes[j].Name
	})
}
