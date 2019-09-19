package resolvers

import (
	"os/exec"

	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
)

type link struct {
	linker string
	flags  []string
}

// NewLink returns a node.Resolver that links a bunch of object files into
// an executable.
func NewLink(linker string, flags []string) node.Resolver {
	return &link{
		linker: linker,
		flags:  flags,
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
	cmd.Args = append(cmd.Args, l.flags...)

	log.Debugf("linker: running %s", cmd.Args)
	o, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(err, string(o))
	}

	return nil
}
