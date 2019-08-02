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

type Resolver interface {
	Resolve(n *Node0) error
}

type Node0 interface {
	Name() string
	Resolve() error

	// Get metadata. A nil return value means the key does not
	// have a corresponding value.
	Get(string) []string
	// Set metadata. A nil second parameter means the key does
	// not have a corresponding value.
	Set(string, []string)
}
