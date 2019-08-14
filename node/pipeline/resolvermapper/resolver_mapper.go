// Package resolvermapper provides a type that can map from a string to a
// node.Resolver.
package resolvermapper

import (
	"fmt"
	"net/http"

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
	ctx *pipeline.Ctx,
	projectDir string,
	n *node.Node,
	name string,
	config map[string]interface{},
) error {
	root := ctx.KV[pipeline.CtxRoot]
	compilerC := ctx.KV[pipeline.CtxCompilerC]
	compilerCC := ctx.KV[pipeline.CtxCompilerCC]
	archiver := ctx.KV[pipeline.CtxArchiver]
	linker := ctx.KV[pipeline.CtxLinker]

	var r node.Resolver
	var err error
	switch name {
	case "compileC":
		r = resolvers.NewCompile(root, compilerC, []string{root})
	case "compileCC":
		r = resolvers.NewCompile(root, compilerCC, []string{root})
	case "archive":
		r = resolvers.NewArchive(root, archiver)
	case "link":
		r = resolvers.NewLink(root, linker)
	case "unzip":
		r = resolvers.NewUnzip(projectDir)
	case "download":
		r, err = createDownload(config)
		if err != nil {
			err = errors.Wrap(err, "create download")
		}
	case "":
		r = nil
	default:
		err = fmt.Errorf("unknown resolver: %s", name)
	}

	if err != nil {
		return err
	}

	n.Resolver = r

	return nil
}

func createDownload(
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
