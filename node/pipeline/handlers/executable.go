package handlers

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/pipeline"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type executable struct {
	s       Store
	rf      ResolverFactory
	project string
	target  string
}

// NewExecutable creates a pipeline.Handler that will add object file node.Node's
// to the node.Node list based on an executable target.
func NewExecutable(
	s Store,
	rf ResolverFactory,
	project string,
	target string,
) pipeline.Handler {
	return &executable{
		s:       s,
		rf:      rf,
		project: project,
		target:  target,
	}
}

func (e *executable) Handle(ctx *pipeline.Ctx) error {
	if filepath.Ext(e.target) != "" {
		return nil
	}

	var dN *node.Node
	sourceCN := node.Find(e.target+".c", ctx.Nodes)
	sourceCCN := node.Find(e.target+".cc", ctx.Nodes)
	if sourceCN != nil && sourceCCN != nil {
		return fmt.Errorf(
			"ambiguous executable %s (%s or %s)",
			e.target,
			sourceCN.Name,
			sourceCCN.Name,
		)
	} else if sourceCN != nil {
		dN = sourceCN
	} else if sourceCCN != nil {
		dN = sourceCCN
	} else {
		return fmt.Errorf("unknown source for executable %s", e.target)
	}

	objectNodes := make([]*node.Node, 0)
	var err error
	objectNodes, err = e.collectObjects(ctx, dN, objectNodes)
	if err != nil {
		return errors.Wrap(err, "collect objects")
	}

	targetN := node.New(e.target)
	for _, objectN := range objectNodes {
		ctx.Nodes = append(ctx.Nodes, objectN)
		targetN.Dependency(objectN)
	}
	name := "link"
	config := make(map[string]interface{})
	targetN.Resolver, err = e.rf.NewResolver(name, config)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("new %s resolver", name))
	}
	ctx.Nodes = append(ctx.Nodes, targetN)

	symlinkN, err := symlinkNFromN(e.rf, targetN, e.target)
	if err != nil {
		return errors.Wrap(err, "symlink from n")
	}
	ctx.Nodes = append(ctx.Nodes, symlinkN)

	return nil
}

func (e *executable) Name() string { return "executable" }

func (e *executable) collectObjects(
	ctx *pipeline.Ctx,
	sourceN *node.Node,
	objectNodes []*node.Node,
) ([]*node.Node, error) {
	logrus.Debugf("collect objects from %s", sourceN.Name)

	objectN, err := objectNFromSourceN(e.s, e.rf, e.project, sourceN)
	if err != nil {
		return nil, errors.Wrap(err, "object from source")
	}

	if node.Find(objectN.Name, objectNodes) != nil {
		return objectNodes, nil
	}
	objectNodes = append(objectNodes, objectN)

	for _, dN := range sourceN.Dependencies {
		ext := filepath.Ext(sourceN.Name)
		source := strings.ReplaceAll(dN.Name, ".h", ext)
		if source == sourceN.Name {
			continue
		}

		sourceN := node.Find(source, ctx.Nodes)
		logrus.Debugf("dependency %s, source %s, found %s", dN, source, sourceN)
		if sourceN != nil {
			objectNodes, err = e.collectObjects(ctx, sourceN, objectNodes)
			if err != nil {
				return nil, err
			}
		}
	}

	return objectNodes, nil
}
