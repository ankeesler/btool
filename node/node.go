// Package node provides a project-wide abstraction of a build-able unit.
//
// Each Node is a file on disk. It depends on other Node's.
//
// Each Node has a Resolver which brings it into existence on disk.
//
// This allows for a very simple "build" algorithm to bring a Node into
// existence:
//   for dependency in node.dependencies:
//     resolve(dependency)
//   resolve(node)
package node

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Resolver

// A Resolver interface is something that brings a node into existence on disk.
//
// For example, it might write out a file that a node is representing.
type Resolver interface {
	Resolve(*Node) error
}

// Node is a build-able unit. Its Name should refer to an actual file on disk.
//
// It has a list of Dependencies of other Node's. It also has a Resolver that
// is in charge of bringing it into existence on disk.
type Node struct {
	Name         string
	Dependencies []*Node
	Resolver
}

// New creates a new Node with a Name and an empty Dependencies list.
func New(name string) *Node {
	return &Node{
		Name:         name,
		Dependencies: make([]*Node, 0),
	}
}

// Dependency adds a list of Nodes to a Node's Dependencies list.
func (n *Node) Dependency(d ...*Node) *Node {
	n.Dependencies = append(n.Dependencies, d...)
	return n
}

// String returns a human-readable representation of a Node.
func (n *Node) String() string {
	//return fmt.Sprintf("%s:%s", n.Name, n.Dependencies)
	return n.Name
}
