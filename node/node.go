// Package node provides a project-wide abstraction of a source unit.
package node

type Resolver interface {
	Resolve(*Node) error
}

type Node struct {
	Name         string
	Dependencies []*Node
	Metadata     map[string][]string

	Resolver
}

func New(name string) *Node {
	return &Node{
		Name:         name,
		Dependencies: make([]*Node, 0),
		Metadata:     make(map[string][]string),
	}
}

func (n *Node) Dependency(d ...*Node) *Node {
	n.Dependencies = append(n.Dependencies, d...)
	return n
}

func (n *Node) String() string {
	//return fmt.Sprintf("%s:%s", n.Name, n.Dependencies)
	return n.Name
}
