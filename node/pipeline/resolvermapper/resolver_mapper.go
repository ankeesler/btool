// Package resolvermapper provides a type that can map from a string to a
// node.Resolver.
package resolvermapper

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/pipeline"
	"github.com/ankeesler/btool/node/resolvers"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// ResolverMapper can return a node.Resolver given a string and some config.
type ResolverMapper struct {
	ctx *pipeline.Ctx
}

// New creates a new ResolverMapper.
func New(ctx *pipeline.Ctx) *ResolverMapper {
	return &ResolverMapper{
		ctx: ctx,
	}
}

// Map returns a node.Resolver given an name and a config map. If no know
// Resolver can be created from the name and the config map, an error should
// be returned.
func (rm *ResolverMapper) Map(
	name string,
	config map[string]interface{},
) (node.Resolver, error) {
	root := rm.ctx.KV[pipeline.CtxRoot]
	cache := rm.ctx.KV[pipeline.CtxCache]
	compilerC := rm.ctx.KV[pipeline.CtxCompilerC]
	compilerCC := rm.ctx.KV[pipeline.CtxCompilerCC]
	archiver := rm.ctx.KV[pipeline.CtxArchiver]
	linker := rm.ctx.KV[pipeline.CtxLinker]
	switch name {
	case "compileC":
		return resolvers.NewCompile(root, compilerC, []string{root}), nil
	case "compileCC":
		return resolvers.NewCompile(root, compilerCC, []string{root}), nil
	case "archive":
		return resolvers.NewArchive(root, archiver), nil
	case "link":
		return resolvers.NewLink(root, linker), nil
	case "unzip":
		r, err := createUnzip(config, cache)
		if err != nil {
			return nil, errors.Wrap(err, "create download")
		}
		return r, nil
	case "download":
		r, err := createDownload(config, cache)
		if err != nil {
			return nil, errors.Wrap(err, "create download")
		}
		return r, nil
	default:
		return nil, fmt.Errorf("unknown resolver: %s", name)
	}
}

func createUnzip(
	config map[string]interface{},
	cache string,
) (node.Resolver, error) {
	outputDir := filepath.Join(cache, "unzip")
	return resolvers.NewUnzip(outputDir), nil
}

func createDownload(
	config map[string]interface{},
	cache string,
) (node.Resolver, error) {
	c := &http.Client{}

	cfg := struct {
		URL    string
		SHA256 string
	}{}
	if err := mapstructure.Decode(config, &cfg); err != nil {
		return nil, errors.Wrap(err, "decode")
	}

	outputFile := filepath.Join(cache, "download", cfg.SHA256)

	return resolvers.NewDownload(c, cfg.URL, cfg.SHA256, outputFile), nil
}
