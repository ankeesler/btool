package gaggle

import (
	"fmt"
	"path/filepath"

	"github.com/ankeesler/btool/collector"
	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/registry"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// Collector is a type that can build a node.Node graph using a registry.Gaggle.
type Collector struct {
}

// New creates a new Collector.
func New() *Collector {
	return &Collector{}
}

func (c *Collector) Collect(
	ctx *collector.Ctx,
	gaggle *registry.Gaggle,
	root string,
) error {
	metadata := struct {
		IncludeDirs []string `mapstructure:"includeDirs"`
	}{}
	if err := mapstructure.Decode(gaggle.Metadata, &metadata); err != nil {
		return errors.Wrap(err, "decode")
	}
	log.Debugf("metadata: %+v", metadata)

	for _, i := range metadata.IncludeDirs {
		i = filepath.Join(root, i)
	}

	for _, n := range gaggle.Nodes {
		nN := node.New(filepath.Join(root, n.Name))
		for _, d := range n.Dependencies {
			dName := filepath.Join(root, d)

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

		nodeR, err := c.newResolver(
			ctx,
			n.Resolver,
			metadata.IncludeDirs,
			root,
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

func (c *Collector) newResolver(
	ctx *collector.Ctx,
	registryR registry.Resolver,
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
		nodeR, err = c.createDownload(ctx, config)
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

func (c *Collector) createDownload(
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
