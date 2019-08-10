// Package registry provides btool registry functionality.
//
// A registry is a place where one can get serialized node.Node's.
package registry

// Registry is how one accesses a registry. There is an Index with the
// listing of Node's, and a list of Node's that can be accessed via Nodes().
type Registry interface {
	// Index returns an Index of the Node's in this registry.
	Index() (*Index, error)
	// Nodes returns the Node's in this registry. It returns nil if there are no
	// know nodes with the provided name, and returns an error if there is any
	// other error.
	Nodes(string) ([]*Node, error)
}

// Index describes the files in the registry.
type Index struct {
	Files []IndexFile `yaml:"files"`
}

func newIndex() *Index {
	return &Index{
		Files: make([]IndexFile, 0),
	}
}

// IndexFile describes a single file in the registry.
type IndexFile struct {
	Path   string `yaml:"path"`
	SHA256 string `yaml:"sha256"`
}

// Node is a serialized form of a node.Node.
type Node struct {
	Name         string   `yaml:"name"`
	Dependencies []string `yaml:"dependencies"`
	Resolver     Resolver `yaml:"resolver"`
}

// String returns a human-readable representation of a Node.
func (n *Node) String() string {
	//return fmt.Sprintf("%s:%s", n.Name, n.Dependencies)
	return n.Name
}

// Resolver is a serialized form of a node.Resolver.
type Resolver struct {
	Name   string                 `yaml:"name"`
	Config map[string]interface{} `yaml:"config"`
}

// Error is a generic error representation.
type Error struct {
	Error string `yaml:"error"`
}
