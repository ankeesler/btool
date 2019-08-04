package objecter

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ankeesler/btool/node"
	"github.com/sirupsen/logrus"
)

type Executabler struct {
}

func NewExecutabler() *Executabler {
	return &Executabler{}
}

func (e *Executabler) Handle(cfg *node.Config, nodes []*node.Node) ([]*node.Node, error) {
	if filepath.Ext(cfg.Target) != "" {
		return nodes, nil
	}

	var d *node.Node
	var ext string
	c := node.Find(cfg.Target+".c", nodes)
	cc := node.Find(cfg.Target+".cc", nodes)
	if c != nil && cc != nil {
		return nil, fmt.Errorf(
			"ambiguous executable %s (%s or %s)",
			cfg.Target,
			c.Name,
			cc.Name,
		)
	} else if c != nil {
		d = c
		ext = ".c"
	} else if cc != nil {
		d = cc
		ext = ".cc"
	} else if d == nil {
		return nil, fmt.Errorf("unknown source for executable %s", cfg.Target)
	}

	var comp string
	if c != nil {
		comp = cfg.CCompiler
	} else { // cc != nil
		comp = cfg.CCCompiler
	}

	objects := make([]*node.Node, 0)
	objects = collectObjects(d, nodes, objects, ext, cfg.Root, comp)

	targetN := node.New(cfg.Target)
	targetN.Resolver = &linker{
		link: comp,
		dir:  cfg.Root,
	}
	for _, object := range objects {
		nodes = append(nodes, object)
		targetN.Dependency(object)
	}
	nodes = append(nodes, targetN)

	return nodes, nil
}

func collectObjects(
	n *node.Node,
	nodes, objects []*node.Node,
	ext, root, comp string,
) []*node.Node {
	logrus.Debugf("collect objects from %s", n.Name)

	objects = append(objects, objectFromSource(n, ext, root, comp))

	for _, d := range n.Dependencies {
		source := strings.ReplaceAll(d.Name, ".h", ext)
		if source == n.Name {
			continue
		}

		sourceN := node.Find(source, nodes)

		logrus.Debugf("dependency %s, source %s, found %s", d, source, sourceN)
		if sourceN != nil {
			objects = collectObjects(sourceN, nodes, objects, ext, root, comp)
		}
	}

	return objects
}

func objectFromSource(source *node.Node, ext, root, comp string) *node.Node {
	object := strings.ReplaceAll(source.Name, ext, ".o")

	logrus.Debugf("adding %s -> %s", object, source.Name)
	n := node.New(object).Dependency(source)
	n.Resolver = &compiler{
		comp:     comp,
		source:   source.Name,
		includes: []string{root},

		dir: root,
	}

	return n
}
