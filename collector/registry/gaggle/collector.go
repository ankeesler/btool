package gaggle

import (
	"fmt"
	"path/filepath"

	"github.com/ankeesler/btool/collector"
	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Gaggler

// Gaggler is an object that can retrieve registry.Gaggle's from somewhere.
type Gaggler interface {
	Gaggle() *registrypkg.Gaggle
	Root() string
}

// Registry is a type that can build a node.Node graph using a
// registrypkg.Gaggle.
type Registry struct {
	g Gaggler
}

// New creates a new Registry with a Gaggler and a root directory.
func New(g Gaggler) *Registry {
	return &Registry{
		g: g,
	}
}

func (r *Registry) Collect(ctx *collector.Ctx, n *node.Node) error {
	gaggle := r.g.Gaggle()

	metadata := struct {
		IncludeDirs []string `mapstructure:"includeDirs"`
	}{}
	if err := mapstructure.Decode(gaggle.Metadata, &metadata); err != nil {
		return errors.Wrap(err, "decode")
	}
	log.Debugf("metadata: %+v", metadata)

	for _, i := range metadata.IncludeDirs {
		i = filepath.Join(r.g.Root(), i)
	}

	for _, n := range gaggle.Nodes {
		nN := node.New(filepath.Join(r.g.Root(), n.Name))
		for _, d := range n.Dependencies {
			dName := filepath.Join(r.g.Root(), d)

			var dN *node.Node
			if d == "$this" {
				// TODO: test me.
				dN = node.New("")
			} else {
				dN = ctx.NS.Find(dName)
			}

			if dN == nil {
				return fmt.Errorf("cannot find dependency %s/%s of %s", d, dName, n)
			}
			nN.Dependency(dN)
		}

		nodeR, err := r.newResolver(
			ctx,
			n.Resolver,
			metadata.IncludeDirs,
			r.g.Root(),
		)
		if err != nil {
			return errors.Wrap(err, "new resolver")
		}
		nN.Resolver = nodeR

		log.Debugf("decoded %s to %s", n, nN)
		ctx.NS.Add(nN)
	}

	return nil
}

func (r *Registry) newResolver(
	ctx *collector.Ctx,
	registryR registrypkg.Resolver,
	includeDirs []string,
	root string,
) (node.Resolver, error) {
	name := registryR.Name
	config := registryR.Config

	var nodeR node.Resolver
	var err error
	switch name {
	case "compileC":
		nodeR = ctx.RF.NewCompileC(includeDirs)
	case "compileCC":
		nodeR = ctx.RF.NewCompileCC(includeDirs)
	case "archive":
		nodeR = ctx.RF.NewArchive()
	case "linkC":
		nodeR = ctx.RF.NewLinkC()
	case "linkCC":
		nodeR = ctx.RF.NewLinkCC()
	case "symlink":
		nodeR = ctx.RF.NewSymlink()
	case "unzip":
		nodeR = ctx.RF.NewUnzip(root)
	case "download":
		nodeR, err = r.createDownload(ctx, config)
		if err != nil {
			err = errors.Wrap(err, "create download")
		}
	case "":
		nodeR = nil
	default:
		err = fmt.Errorf("unknown resolver: %s", name)
	}

	return nodeR, err
}

func (r *Registry) createDownload(
	ctx *collector.Ctx,
	config map[string]interface{},
) (node.Resolver, error) {
	cfg := struct {
		URL    string
		SHA256 string
	}{}
	if err := mapstructure.Decode(config, &cfg); err != nil {
		return nil, errors.Wrap(err, "decode")
	}

	return ctx.RF.NewDownload(cfg.URL, cfg.SHA256), nil
}
