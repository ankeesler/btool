// Package registry provideds an abstraction of a collection of serialized
// node.Node's.
package registry

import "github.com/ankeesler/btool/node"

type Node struct {
	Name         string
	Dependencies []string
	Labels       map[string]string
}

type Client interface {
	List() ([]*node.Node, error)
	Get(string) (*node.Node, error)
}
