package resolvers

import (
	"os/exec"

	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type archive struct {
	archiver string
}

// NewArchive returns a node.Resolver that runs the archiver in order to generate
// a static library.
func NewArchive(archiver string) node.Resolver {
	return &archive{
		archiver: archiver,
	}
}

func (a *archive) Resolve(n *node.Node) error {
	logrus.Debugf("archiver: resolve %s/%s", n, n.Dependencies)

	cmd := exec.Command(
		a.archiver,
		"rcs",
		n.Name,
	)
	for _, dN := range n.Dependencies {
		cmd.Args = append(cmd.Args, dN.Name)
	}

	logrus.Debugf("archiver: running %s from %s", cmd.Args, cmd.Dir)
	o, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(err, string(o))
	}

	return nil
}
