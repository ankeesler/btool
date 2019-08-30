package collector

import (
	"github.com/ankeesler/btool/node"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . ResolverFactory

// ResolverFactory can create node.Resolver's.
type ResolverFactory interface {
	NewCompileC(includeDirs []string) node.Resolver
	NewCompileCC(includeDirs []string) node.Resolver
	NewArchive() node.Resolver
	NewLinkC() node.Resolver
	NewLinkCC() node.Resolver
	NewSymlink() node.Resolver

	NewDownload(url, sha256 string) node.Resolver
	NewUnzip(outputDir string) node.Resolver
}

// Ctx passes shared objects through to Collectini's.
type Ctx struct {
	NS *NodeStore
	RF ResolverFactory

	includePaths []string
	libraries    map[string][]*node.Node
}

// NewCtx creates a new Ctx with a NodeStore and a ResolverFactory.
func NewCtx(ns *NodeStore, rf ResolverFactory) *Ctx {
	return &Ctx{
		NS: ns,
		RF: rf,

		includePaths: make([]string, 0),
		libraries:    make(map[string][]*node.Node),
	}
}

func (c *Ctx) AddIncludePath(includePath string) {
	c.includePaths = append(c.includePaths, includePath)
}

func (c *Ctx) IncludePaths() []string {
	return c.includePaths
}

func (c *Ctx) AddLibrary(include string, libN *node.Node) {
	libraries, ok := c.libraries[include]
	if !ok {
		libraries = make([]*node.Node, 0)
	}
	c.libraries[include] = append(libraries, libN)
}

func (c *Ctx) Libraries(include string) ([]*node.Node, bool) {
	libraries, ok := c.libraries[include]
	return libraries, ok
}
