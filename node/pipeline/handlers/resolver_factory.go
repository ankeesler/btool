package handlers

import "github.com/ankeesler/btool/node"

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . ResolverFactory

// ResolverFactory can create new node.Resolver's by a name and a config.
type ResolverFactory interface {
	NewResolver(
		name string,
		config map[string]interface{},
	) (node.Resolver, error)
}
