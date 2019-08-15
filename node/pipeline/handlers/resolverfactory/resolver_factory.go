// Package resolverfactory provides a factory type for creating node.Resolver's.
package resolverfactory

import (
	"net/http"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/resolvers"
)

// ResolverFactory is a factory type that can create node.Resolver's.
type ResolverFactory struct {
	compilerC, compilerCC, archiver, linker string
}

// New creates a new ResolverFactory.
func New(
	compilerC, compilerCC, archiver, linker string,
) *ResolverFactory {
	return &ResolverFactory{
		compilerC:  compilerC,
		compilerCC: compilerCC,
		archiver:   archiver,
		linker:     linker,
	}
}

// NewCompileC returns a node.Resolver that compiles a C source into an object.
func (rf *ResolverFactory) NewCompileC(includePaths []string) node.Resolver {
	return resolvers.NewCompile(rf.compilerC, includePaths)
}

// NewCompileCC returns a node.Resolver that compiles a C++ source into an
// object.
func (rf *ResolverFactory) NewCompileCC(includePaths []string) node.Resolver {
	return resolvers.NewCompile(rf.compilerCC, includePaths)
}

// NewArchive returns a node.Resolver that archives multiple objects into a
// library.
func (rf *ResolverFactory) NewArchive() node.Resolver {
	return resolvers.NewArchive(rf.archiver)
}

// NewLink returns a node.Resolver that links multiple objects into an
// executable.
func (rf *ResolverFactory) NewLink() node.Resolver {
	return resolvers.NewLink(rf.linker)
}

// NewSymlink returns a node.Resolver that symlinks a node to another.
func (rf *ResolverFactory) NewSymlink() node.Resolver {
	return resolvers.NewSymlink()
}

// NewDownload returns a node.Resolver that downloads an HTTP/HTTPS URL onto
// disk.
func (rf *ResolverFactory) NewDownload(url, sha256 string) node.Resolver {
	return resolvers.NewDownload(&http.Client{}, url, sha256)
}

// NewUnzip returns a node.Resolver that unzips a zip archive into a directory.
func (rf *ResolverFactory) NewUnzip(outputDir string) node.Resolver {
	return resolvers.NewUnzip(outputDir)
}
