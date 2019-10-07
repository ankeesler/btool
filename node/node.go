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
//
// A Node also has a simple map of Labels to which clients can attach metadata.
type Node struct {
	Name         string
	Dependencies []*Node
	Labels       map[string]interface{}
	Resolver
}

// New creates a new Node with a Name and an empty Dependencies list.
func New(name string) *Node {
	return &Node{
		Name:         name,
		Dependencies: make([]*Node, 0),
		Labels:       make(map[string]interface{}),
	}
}

// Dependency adds a list of Nodes to a Node's Dependencies list.
// If a Dependency already exists (i.e., if there is already a Node in the
// Dependencies list with the same Name as the new Dependency), it will be
// replaced.
func (n *Node) Dependency(d ...*Node) *Node {
	for _, dd := range d {
		n.replaceDependency(dd)
	}
	return n
}

func (n *Node) replaceDependency(d *Node) {
	for i := range n.Dependencies {
		if n.Dependencies[i].Name == d.Name {
			n.Dependencies[i] = d
			return
		}
	}
	n.Dependencies = append(n.Dependencies, d)
}

// Label adds a Label to a Node. It returns this Node so calls can be strung
// together.
func (n *Node) Label(k string, v interface{}) *Node {
	n.Labels[k] = v
	return n
}

// String returns a human-readable representation of a Node.
func (n *Node) String() string {
	return n.Name
}
