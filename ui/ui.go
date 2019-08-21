// Package ui provides btool command line user interface with pretty printing.
package ui

import (
	"fmt"
	"strings"

	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
)

// UI is an object that can provide pretty btool command line printing.
type UI struct {
	quiet bool
}

// New creates a new UI.
func New(quiet bool) *UI {
	return &UI{
		quiet: quiet,
	}
}

func (ui *UI) OnResolve(n *node.Node, current bool) {
	if ui.quiet {
		return
	}

	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("resolving %s", n.Name))
	if current {
		b.WriteString(" (up to date)")
	}
	log.Infof(b.String())
}

func (ui *UI) OnAdd(n *node.Node) {
	if ui.quiet {
		return
	}

	log.Infof("adding " + n.Name)
}

func (ui *UI) OnClean(n *node.Node) {
	if ui.quiet {
		return
	}

	log.Infof("cleaning " + n.Name)
}
