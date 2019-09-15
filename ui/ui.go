// Package ui provides btool command line user interface with pretty printing.
package ui

import (
	"fmt"
	"strings"

	collector "github.com/ankeesler/btool/collector0"
	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
)

// UI is an object that can provide pretty btool command line printing.
type UI struct {
	quiet bool

	added map[*node.Node]bool
}

// New creates a new UI.
func New(quiet bool) *UI {
	return &UI{
		quiet: quiet,

		added: make(map[*node.Node]bool),
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

func (ui *UI) Consume(s collector.Store, n *node.Node) error {
	if ui.quiet {
		return nil
	}

	if _, ok := ui.added[n]; !ok {
		log.Infof("adding " + n.Name)
		ui.added[n] = true
	}

	return nil
}

func (ui *UI) OnClean(n *node.Node) {
	if ui.quiet {
		return
	}

	log.Infof("cleaning " + n.Name)
}
