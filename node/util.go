package node

func Find(target string, nodes []*Node) *Node {
	for _, n := range nodes {
		if n.Name == target {
			return n
		}
	}
	return nil
}
