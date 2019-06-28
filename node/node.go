// Package node provides a project-wide abstraction of a source unit.
package node

type Node struct {
	// Unique name.
	Name string

	// On-disk .c/.cc and .h files.
	Sources []string
	Headers []string

	Dependencies []*Node

	// -I include paths needed in compiler invocation.
	IncludePaths []string

	// Object archives for the Sources. On-disk files.
	Objects []string
}

func (n *Node) String() string {
	//return fmt.Sprintf("%s:%s", n.Name, n.Dependencies)
	return n.Name
}
