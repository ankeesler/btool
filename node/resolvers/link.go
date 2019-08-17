package resolvers

import (
	"os/exec"

	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
)

type link struct {
	linker string
}

// NewLink returns a node.Resolver that links a bunch of object files into
// an executable.
func NewLink(linker string) node.Resolver {
	return &link{
		linker: linker,
	}
}

func (l *link) Resolve(n *node.Node) error {
	cmd := exec.Command(
		l.linker,
		"-o",
		n.Name,
	)
	for _, d := range n.Dependencies {
		cmd.Args = append(cmd.Args, d.Name)
	}

	log.Debugf("linker: running %s from %s", cmd.Args, cmd.Dir)
	o, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(err, string(o))
	}

	return nil
}
