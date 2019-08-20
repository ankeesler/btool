package handlers

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/pipeline"
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

func (e *executable) Handle(ctx pipeline.Ctx) error {
	if filepath.Ext(e.target) != "" {
		return nil
	}

	var dN *node.Node
	sourceCN := node.Find(e.target+".c", ctx.All())
	sourceCCN := node.Find(e.target+".cc", ctx.All())
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
	objectNodes = e.collectObjects(ctx, dN, objectNodes)

	targetN := node.New(e.target)
	for _, objectN := range objectNodes {
		ctx.Add(objectN)
		targetN.Dependency(objectN)
	}
	if sourceCCN != nil {
		targetN.Resolver = e.rf.NewLinkCC()
	} else {
		targetN.Resolver = e.rf.NewLinkC()
	}
	ctx.Add(targetN)

	return nil
}

func (e *executable) String() string { return "executable" }

func (e *executable) collectObjects(
	ctx pipeline.Ctx,
	sourceN *node.Node,
	objectNodes []*node.Node,
) []*node.Node {
	log.Debugf("collect objects from %s", sourceN.Name)

	objectN := objectNFromSourceN(e.s, e.rf, e.project, sourceN)
	if node.Find(objectN.Name, objectNodes) != nil {
		return objectNodes
	}
	objectNodes = append(objectNodes, objectN)

	for _, dN := range sourceN.Dependencies {
		ext := filepath.Ext(sourceN.Name)
		source := strings.ReplaceAll(dN.Name, ".h", ext)
		if source == sourceN.Name {
			continue
		}

		sourceN := ctx.Find(source)
		log.Debugf("dependency %s, source %s, found %s", dN, source, sourceN)
		if sourceN != nil {
			objectNodes = e.collectObjects(ctx, sourceN, objectNodes)
		}
	}

	return objectNodes
}
