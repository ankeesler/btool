package graph

type Node struct {
	Name string
}

func (n *Node) String() string { return n.Name }
