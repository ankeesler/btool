package registry

import (
	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . NodeMapper

// NodeMapper returns a node.Node given a string.
type NodeMapper interface {
	// Map returns a node.Node given a string. If there is no known node.Node for
	// the provided string, then an error should be returned.
	Map(string) (*node.Node, error)
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . ResolverMapper

// ResolverMapper returns a node.Resolver given a name and a config map.
type ResolverMapper interface {
	// Map returns a node.Resolver given an name and a config map. If no know
	// Resolver can be created from the name and the config map, an error should
	// be returned.
	Map(string, map[string]interface{}) (node.Resolver, error)
}

// Decoder can turn a Node into a node.Node.
type Decoder struct {
	nodeMapper     NodeMapper
	resolverMapper ResolverMapper
}

// NewDecoder creates a new Decoder with a NodeMapper and a ResolverMapper.
func NewDecoder(nodeMapper NodeMapper, resolverMapper ResolverMapper) *Decoder {
	return &Decoder{
		nodeMapper:     nodeMapper,
		resolverMapper: resolverMapper,
	}
}

// Decode turns a Node into a node.Node.
func (d *Decoder) Decode(n *Node) (*node.Node, error) {
	nN := node.New(n.Name)
	for _, dependency := range n.Dependencies {
		dependencyN, err := d.nodeMapper.Map(dependency)
		if err != nil {
			return nil, errors.Wrap(err, "map node")
		}
		nN.Dependency(dependencyN)
	}

	r, err := d.resolverMapper.Map(n.Resolver.Name, n.Resolver.Config)
	if err != nil {
		return nil, errors.Wrap(err, "map resolver")
	}
	nN.Resolver = r

	return nN, nil
}
