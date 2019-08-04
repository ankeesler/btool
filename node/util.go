package node

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
