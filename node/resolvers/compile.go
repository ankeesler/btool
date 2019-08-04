package resolvers

import (
	"fmt"
	"os/exec"

	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type compile struct {
	dir      string
	compiler string
	includes []string
}

// NewCompile returns a node.Resolver that runs a compiler in order to generate
// an object file.
func NewCompile(dir, compiler string, includes []string) node.Resolver {
	return &compile{
		dir:      dir,
		compiler: compiler,
		includes: includes,
	}
}

func (c *compile) Resolve(n *node.Node) error {
	if len(n.Dependencies) != 1 {
		return fmt.Errorf("expected %d dependencies, got %d", 1, len(n.Dependencies))
	}

	cmd := exec.Command(
		c.compiler,
		"-o",
		n.Name,
		"-c",
		n.Dependencies[0].Name,
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
