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

type object struct {
	s       Store
	rf      ResolverFactory
	project string
	target  string
}

// NewObject creates a pipeline.Handler that will add object file node.Node's to
// the node.Node list based on an object target.
func NewObject(
	s Store,
	rf ResolverFactory,
	project string,
	target string,
) pipeline.Handler {
	return &object{
		s:       s,
		rf:      rf,
		project: project,
		target:  target,
	}
}

func (o *object) Handle(ctx *pipeline.Ctx) error {
	if !strings.HasSuffix(o.target, ".o") {
		return nil
	}

	var dN *node.Node
	sourceCN := node.Find(strings.ReplaceAll(o.target, ".o", ".c"), ctx.Nodes)
	sourceCCN := node.Find(strings.ReplaceAll(o.target, ".o", ".cc"), ctx.Nodes)
	if sourceCN != nil && sourceCCN != nil {
		return fmt.Errorf(
			"ambiguous object %s (%s or %s)",
			o.target,
			sourceCN.Name,
			sourceCCN.Name,
		)
	} else if sourceCN != nil {
		dN = sourceCN
	} else if sourceCCN != nil {
		dN = sourceCCN
	} else {
		return fmt.Errorf("unknown source for object %s", o.target)
	}

	objectN, err := objectNFromSourceN(o.s, o.rf, o.project, dN)
	if err != nil {
		return errors.Wrap(err, "object from source")
	}
	ctx.Nodes = append(ctx.Nodes, objectN)

	symlinkN, err := symlinkNFromN(o.rf, objectN, o.target)
	if err != nil {
		return errors.Wrap(err, "symlink from n")
	}
	ctx.Nodes = append(ctx.Nodes, symlinkN)

	return nil
}

func (o *object) Name() string { return "object" }

func objectNFromSourceN(
	s Store,
	rf ResolverFactory,
	project string,
	sourceN *node.Node,
) (*node.Node, error) {
	ext := filepath.Ext(sourceN.Name)

	name := "compile" + ext
	config := map[string]interface{}{
		"includePaths": []string{s.ProjectDir(project)},
	}
	r, err := rf.NewResolver(name, config)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("new %s resolver", name))
	}

	object := strings.ReplaceAll(sourceN.Name, ext, ".o")
	logrus.Debugf(
		"adding %s -> %s",
		object,
		sourceN.Name,
	)
	objectN := node.New(object).Dependency(sourceN)
	objectN.Resolver = r
	return objectN, nil
}

func symlinkNFromN(
	rf ResolverFactory,
	n *node.Node,
	target string,
) (*node.Node, error) {
	name := "symlink"
	config := make(map[string]interface{})
	r, err := rf.NewResolver("symlink", config)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("new %s resolver", name))
	}

	symlinkN := node.New(target).Dependency(n)
	symlinkN.Resolver = r

	return symlinkN, nil
}
