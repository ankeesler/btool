package node

import (
	"sort"
)

// Find is a utility function that searches for a Node with the provided name
// in a list of Node's. It will return nil if no such Node exists.
func Find(name string, nodes []*Node) *Node {
	for _, n := range nodes {
		if n.Name == name {
			return n
		}
	}
	return nil
}

// SortaAlpha sorts the provided node.Node's and dependencies in alphanumeric
// order.
func SortAlpha(nodes []*Node) {
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
