package collector

import (
	"path/filepath"

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

	includePaths map[string]string
}

// NewCtx creates a new Ctx with a NodeStore and a ResolverFactory.
func NewCtx(ns *NodeStore, rf ResolverFactory) *Ctx {
	return &Ctx{
		NS: ns,
		RF: rf,

		includePaths: make(map[string]string),
	}
}

// SetIncludePath will mark a node (e.g., gtest/gtest.h) as needing an
// include directory on compiler invocation. After setting include paths with
// this function, IncludePath() can be called to collect all of the necessary
// include paths for a compiler invocation.
func (c *Ctx) SetIncludePath(n *node.Node, path string) {
	c.includePaths[n.Name] = path
}

// IncludePath returns the include path for an include, or "" if there is no
// know include path.
func (c *Ctx) IncludePath(include string) string {
	for _, includePath := range c.includePaths {
		if c.NS.Find(filepath.Join(includePath, include)) != nil {
			return includePath
		}
	}
	return ""
}

// IncludePaths collects all of the necessary include paths for a compiler
// invocation.
func (c *Ctx) IncludePaths(n *node.Node) []string {
	includePaths := make([]string, 0)
	node.Visit(n, func(nV *node.Node) error {
		if includePath, ok := c.includePaths[nV.Name]; ok {
			if !contains(includePaths, includePath) {
				includePaths = append(includePaths, includePath)
			}
		}
		return nil
	})
	return includePaths
}

func contains(ss []string, s string) bool {
	for i := range ss {
		if ss[i] == s {
			return true
		}
	}
	return false
}
