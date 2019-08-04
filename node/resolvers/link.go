package resolvers

import (
	"os/exec"

	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type link struct {
	dir    string
	linker string
}

// NewLink returns a node.Resolver that links a bunch of object files into
// an executable.
func NewLink(dir, linker string) node.Resolver {
	return &link{
		dir:    dir,
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
	cmd.Dir = l.dir

	logrus.Debugf("linker: running %s from %s", cmd.Args, cmd.Dir)
	o, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(err, string(o))
	}

	return nil
}
