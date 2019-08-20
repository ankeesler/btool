package btool

import (
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/builder"
	"github.com/ankeesler/btool/node/builder/currenter"
	"github.com/ankeesler/btool/ui"
	"github.com/pkg/errors"
)

// Build will run create all node.Node's and dependencies in the node.Node graph.
func Build(targetN *node.Node, ui *ui.UI) error {
	c := currenter.New()
	b := builder.New(c, ui)
	if err := b.Build(targetN); err != nil {
		return errors.Wrap(err, "build")
	}

	return nil
}
