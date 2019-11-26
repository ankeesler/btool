package resolvers

import (
	"os/exec"
	"path/filepath"

	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
)

type untar struct {
}

// NewUntar returns a node.Resolver that untars a dependency.
func NewUntar() node.Resolver {
	return &untar{}
}

func (u *untar) Resolve(n *node.Node) error {
	if len(n.Dependencies) == 0 {
		return errors.New("untar resolver target " + n.Name + " must have at least one dependency")
	}

	outputDir := filepath.Dir(n.Dependencies[0].Name)

	cmd := exec.Command("tar")
	cmd.Args = append(cmd.Args, "xzf")
	cmd.Args = append(cmd.Args, n.Dependencies[0].Name)
	cmd.Args = append(cmd.Args, "-C")
	cmd.Args = append(cmd.Args, outputDir)

	log.Debugf("untar: running %s", cmd.Args)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(err, string(output))
	}

	return nil
}
