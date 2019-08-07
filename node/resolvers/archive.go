package resolvers

import (
	"os/exec"

	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type archive struct {
	dir      string
	archiver string
}

// NewArchive returns a node.Resolver that runs the archiver in order to generate
// a static library.
func NewArchive(dir, archiver string) node.Resolver {
	return &archive{
		dir:      dir,
		archiver: archiver,
	}
}

func (a *archive) Resolve(n *node.Node) error {
	cmd := exec.Command(
		a.archiver,
		"rcs",
		n.Name,
	)
	for _, dN := range n.Dependencies {
		cmd.Args = append(cmd.Args, dN.Name)
	}
	cmd.Dir = a.dir

	logrus.Debugf("archiver: running %s from %s", cmd.Args, cmd.Dir)
	o, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(err, string(o))
	}

	return nil
}
