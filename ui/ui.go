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
}

// New creates a new UI.
func New() *UI {
	return &UI{}
}

func (ui *UI) OnResolve(name string, current bool) {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("resolving %s", name))
	if current {
		b.WriteString(" (up to date)")
	}
	log.Infof(b.String())
}

func (ui *UI) OnAdd(n *node.Node) {
	log.Infof("adding " + n.Name)
}
