package btool

import (
	"fmt"
	"strings"

	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
)

type ui int

func newUI() ui {
	return 0
}

func (ui ui) OnResolve(name string, current bool) {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("resolving %s", name))
	if current {
		b.WriteString(" (up to date)")
	}
	log.Infof(b.String())
}

func (ui ui) OnAdd(n *node.Node) {
	log.Infof("adding " + n.Name)
}
