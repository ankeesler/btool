package resolvers

import (
	"fmt"
	"os/exec"

	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
)

type compile struct {
	compiler string
	includes []string
	flags    []string
}

// NewCompile returns a node.Resolver that runs a compiler in order to generate
// an object file.
func NewCompile(compiler string, includes, flags []string) node.Resolver {
	return &compile{
		compiler: compiler,
		includes: includes,
		flags:    flags,
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
	)
	for _, include := range c.includes {
		cmd.Args = append(cmd.Args, "-I"+include)
	}
	cmd.Args = append(cmd.Args, c.flags...)

	log.Debugf("compiler: running %s", cmd.Args)
	o, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(err, string(o))
	}

	return nil
}
