// Package resolverfactory provides a factory type for creating node.Resolver's.
package resolverfactory

import (
	"net/http"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/resolvers"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// These are the known node.Resolver name's that can be provided to
// ResolverFactory.NewResolver(). They should be self explanatory. :)
const (
	nameCompileC  = "compile.c"
	nameCompileCC = "compile.cc"
	nameArchive   = "archive"
	nameLink      = "link"

	nameDownload = "download"
	nameUnzip    = "unzip"
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
	case nameCompileC:
		r, err = rf.createCompile(config, false)
	case nameCompileCC:
		r, err = rf.createCompile(config, true)
	case nameArchive:
		r = resolvers.NewArchive(rf.archiver)
	case nameLink:
		r = resolvers.NewLink(rf.linker)

	case nameDownload:
		r, err = rf.createDownload(config)
	case nameUnzip:
		r, err = rf.createUnzip(config)
	}

	if err != nil {
		return nil, errors.Wrap(err, "create "+name)
	} else if r == nil {
		return nil, errors.New("unknown resolver: " + name)
	}

	return r, nil
}

func (rf *ResolverFactory) createCompile(
	config map[string]interface{},
	cc bool,
) (node.Resolver, error) {
	var compiler string
	if cc {
		compiler = rf.compilerC
	} else {
		compiler = rf.compilerCC
	}

	cfg := struct {
		IncludePaths []string
	}{}
	if err := mapstructure.Decode(config, &cfg); err != nil {
		return nil, errors.Wrap(err, "decode")
	}

	return resolvers.NewCompile(compiler, cfg.IncludePaths), nil
}

func (rf *ResolverFactory) createDownload(
	config map[string]interface{},
) (node.Resolver, error) {
	c := &http.Client{}

	cfg := struct {
		URL    string
		SHA256 string
	}{}
	if err := mapstructure.Decode(config, &cfg); err != nil {
		return nil, errors.Wrap(err, "decode")
	}

	return resolvers.NewDownload(c, cfg.URL, cfg.SHA256), nil
}

func (rf *ResolverFactory) createUnzip(
	config map[string]interface{},
) (node.Resolver, error) {
	cfg := struct {
		OutputDir string
	}{}
	if err := mapstructure.Decode(config, &cfg); err != nil {
		return nil, errors.Wrap(err, "decode")
	}

	return resolvers.NewUnzip(cfg.OutputDir), nil
}
