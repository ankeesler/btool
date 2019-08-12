package node

import "fmt"

// Mapper returns a node.Node given a string. It does this by referencing a list
// of Node's.
type Mapper struct {
	nodes *[]*Node
}

// NewMapper returns a new Mapper with a list of Node's with which to lookup
// names.
func NewMapper(nodes *[]*Node) *Mapper {
	return &Mapper{
		nodes: nodes,
	}
}

// Map returns a Node given a string. If there is no known Node for
// the provided string, then an error should be returned.
func (m *Mapper) Map(name string) (*Node, error) {
	n := Find(name, *m.nodes)
	if n == nil {
		return nil, fmt.Errorf("unknown node: %s", name)
	} else {
		return n, nil
	}
}
