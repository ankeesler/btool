package gaggle

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ankeesler/btool/collector"
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

		for k, v := range n.Labels {
			vBuf := bytes.NewBuffer([]byte{})
			// TODO: this is jank, we should have more of a better interface for this.
			for _, vv := range strings.Split(v, ",") {
				fmt.Fprintf(vBuf, "%s,", filepath.Join(root, vv))
			}
			nN.Label(k, vBuf.String())
		}

		// TODO: is this bad to collect include paths from dependencies first?
		// TODO: this is duplicated code.
		includePaths := make([]string, 0)
		node.Visit(nN, func(vn *node.Node) error {
			// TODO: this shouldn't be hardcoded.
			if ips, ok := vn.Labels["io.btool.cc.includePaths"]; ok {
				// TODO: this is jank, we should have more of a better interface for this.
				for _, ip := range strings.Split(ips, ",") {
					if ip != "" {
						includePaths = append(includePaths, ip)
					}
				}
			}
			return nil
		})

		nodeR, err := c.newResolver(n.Resolver, root, includePaths)
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
		nodeR = c.rf.NewLinkC()
	case "linkCC":
		nodeR = c.rf.NewLinkCC()
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
