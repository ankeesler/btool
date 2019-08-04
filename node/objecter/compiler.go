package objecter

import (
	"os/exec"

	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type compiler struct {
	comp     string
	source   string
	includes []string
	dir      string
}

func (c *compiler) Resolve(n *node.Node) error {
	cmd := exec.Command(
		c.comp,
		"-o",
		n.Name,
		"-c",
		c.source,
		"-Wall",
		"-Werror",
		"-g",
		"-O0",
	)
	for _, include := range c.includes {
		cmd.Args = append(cmd.Args, "-I"+include)
	}
	cmd.Dir = c.dir

	logrus.Debugf("compiler: running %s from %s", cmd.Args, cmd.Dir)
	o, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(err, string(o))
	}

	return nil
}
