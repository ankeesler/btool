package collector

import "github.com/ankeesler/btool/node"

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
}

// NewCtx creates a new Ctx with a NodeStore and a ResolverFactory.
func NewCtx(ns *NodeStore, rf ResolverFactory) *Ctx {
	return &Ctx{
		NS: ns,
		RF: rf,
	}
}
