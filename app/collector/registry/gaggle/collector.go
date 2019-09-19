package gaggle

import (
	"fmt"
	"path/filepath"

	"github.com/ankeesler/btool/app/collector"
	"github.com/ankeesler/btool/app/collector/cc"
	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/registry"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// Collector is a type that can build a node.Node graph using a registry.Gaggle.
type Collector struct {
	rf collector.ResolverFactory
}

// New creates a new Collector.
func New(rf collector.ResolverFactory) *Collector {
	return &Collector{
		rf: rf,
	}
}

func (c *Collector) Collect(
	s collector.Store,
	gaggle *registry.Gaggle,
	root string,
) error {
	for _, n := range gaggle.Nodes {
		nN := node.New(filepath.Join(root, n.Name))

		for _, d := range n.Dependencies {
			dName := filepath.Join(root, d)

			var dN *node.Node
			if d == "$this" {
				// TODO: implement me.
				//dN = node.New("")
				continue
			} else {
				dN = s.Get(dName)
			}

			if dN == nil {
				return fmt.Errorf("cannot find dependency %s of %s", d, n)
			}
			nN.Dependency(dN)
		}

		nN.Labels = n.Labels

		if err := prependRoot(nN, root); err != nil {
			return errors.Wrap(err, "prepend root")
		}

		// TODO: is this bad to reach across sibling packages?
		includePaths, err := cc.CollectIncludePaths(nN)
		if err != nil {
			return errors.Wrap(err, "collect include paths")
		}
		linkFlags, err := cc.CollectLinkFlags(nN)
		if err != nil {
			return errors.Wrap(err, "collect include paths")
		}

		nodeR, err := c.newResolver(n.Resolver, root, includePaths, linkFlags)
		if err != nil {
			return errors.Wrap(err, "new resolver")
		}
		nN.Resolver = nodeR

		log.Debugf("decoded %s to %s", n, nN)
		s.Set(nN)
	}

	return nil
}

func (c *Collector) newResolver(
	registryR registry.Resolver,
	root string,
	includePaths []string,
	linkFlags []string,
) (node.Resolver, error) {
	name := registryR.Name
	config := registryR.Config

	var nodeR node.Resolver
	var err error
	switch name {
	case "compileC":
		nodeR = c.rf.NewCompileC(includePaths)
	case "compileCC":
		nodeR = c.rf.NewCompileCC(includePaths)
	case "archive":
		nodeR = c.rf.NewArchive()
	case "linkC":
		nodeR = c.rf.NewLinkC(linkFlags)
	case "linkCC":
		nodeR = c.rf.NewLinkCC(linkFlags)
	case "symlink":
		nodeR = c.rf.NewSymlink()
	case "unzip":
		nodeR = c.rf.NewUnzip(root)
	case "download":
		nodeR, err = c.createDownload(config)
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
	config map[string]interface{},
) (node.Resolver, error) {
	cfg := struct {
		URL    string
		SHA256 string
	}{}
	if err := mapstructure.Decode(config, &cfg); err != nil {
		return nil, errors.Wrap(err, "decode")
	}

	return c.rf.NewDownload(cfg.URL, cfg.SHA256), nil
}

func prependRoot(n *node.Node, root string) error {
	// TODO: is this bad to reach across sibling packages?
	var labels cc.Labels
	if err := collector.FromLabels(n, &labels); err != nil {
		return errors.Wrap(err, "from labels")
	}

	if labels.IncludePaths == nil {
		labels.IncludePaths = []string{}
	}

	for i := range labels.IncludePaths {
		labels.IncludePaths[i] = filepath.Join(root, labels.IncludePaths[i])
	}

	if labels.Libraries == nil {
		labels.Libraries = []string{}
	}

	for i := range labels.Libraries {
		labels.Libraries[i] = filepath.Join(root, labels.Libraries[i])
	}

	if err := collector.ToLabels(n, &labels); err != nil {
		return errors.Wrap(err, "to labels")
	}

	return nil
}
