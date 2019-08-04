package objecter

import (
	"os/exec"

	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type linker struct {
	link string
	dir  string
}

func (l *linker) Resolve(n *node.Node) error {
	cmd := exec.Command(
		l.link,
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
