// Package resolverfactory provides a factory type for creating node.Resolver's.
package resolverfactory

import (
	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
)

// These are the known node.Resolver name's that can be provided to
// ResolverFactory.NewResolver(). They should be self explanatory. :)
const (
	nameCompile = "compile"
	nameArchive = "archive"
	nameLink    = "link"

	nameDownload = "download"
	nameUnzip    = "unzip"
)

// These are the known config options for each node.Resolver that can be provided
// to ResolverFactory.NewResolver() in the config. They should be self
// explanatory. :)
const (
	configCompileIncludePaths = NameCompile + ".includePaths"
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

// NewResolver creates a new node.Resolver from the provided name and config.
//
// NewResolver will return an error if it cannot create a node.Resolver for the
// provided name and config.
func (rf *ResolverFactory) NewResolver(
	name string,
	config map[string]interface{},
) (node.Resolver, error) {
	var r node.Resolver
	var err error
	switch name {
	}

	if err != nil {
		return errors.Wrap(err, "create "+name)
	} else if r == nil {
		return errors.New("unknown resolver: " + name)
	}

	return r, nil
}
