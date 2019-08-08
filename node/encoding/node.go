package encoding

// Node is a serialized form of a node.Node.
type Node struct {
	Name         string   `yaml:"name"`
	Dependencies []string `yaml:"dependencies"`
	Resolver     Resolver `yaml:"resolver"`
}

// Resolver is a serialized form of a node.Resolver.
type Resolver struct {
	Name   string                 `yaml:"name"`
	Config map[string]interface{} `yaml:"config"`
}
