// Package registry provides btool registry functionality.
//
// A registry is a place where one can get serialized node.Node's.
package registry

// Index describes the files in the registry.
type Index struct {
	Name  string      `yaml:"name"`
	Files []IndexFile `yaml:"files"`
}

func newIndex(name string) *Index {
	return &Index{
		Name:  name,
		Files: make([]IndexFile, 0),
	}
}

// IndexFile describes a single file in the registry.
type IndexFile struct {
	Path   string `yaml:"path"`
	SHA256 string `yaml:"sha256"`
}

// Gaggle is a group of Node's with some Metadata.
type Gaggle struct {
	Metadata map[string]interface{} `yaml:"metadata"`
	Nodes    []*Node                `yaml:"nodes"`
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

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Client

// Registry is an object that can retrieve Gaggle's from a btool registry.
type Registry interface {
	// Index should return the Index associated with this particular
	// Registry. If any error occurs, an error should be returned.
	Index() (*Index, error)
	// Gaggle should return the Gaggle associated with the provided
	// IndexFile.Path. If any error occurs, an error should be returned.
	// If no Gaggle exists for the provided string, then nil, nil should
	// be returned.
	Gaggle(string) (*Gaggle, error)
}
