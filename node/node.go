// Package node provides a project-wide abstraction of a source unit.
package node

type Node struct {
	Name         string
	Sources      []string
	Headers      []string
	Dependencies []*Node
	IncludePaths []string
}

func (n *Node) String() string {
	//return fmt.Sprintf("%s:%s", n.Name, n.Dependencies)
	return n.Name
}
