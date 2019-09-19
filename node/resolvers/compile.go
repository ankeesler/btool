package resolvers

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
)

type compile struct {
	compiler string
	includes []string
}

// NewCompile returns a node.Resolver that runs a compiler in order to generate
// an object file.
func NewCompile(compiler string, includes []string) node.Resolver {
	return &compile{
		compiler: compiler,
		includes: includes,
	}
}

func (c *compile) Resolve(n *node.Node) error {
	if len(n.Dependencies) != 1 {
		return fmt.Errorf("expected %d dependencies, got %d", 1, len(n.Dependencies))
	}

	// TODO: wow what a hack! We should provide these flags another way...
	var std string
	if strings.HasSuffix(c.compiler, "++") {
		std = "c++17"
	} else {
		std = "c17"
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
		"--std="+std,
	)
	for _, include := range c.includes {
		cmd.Args = append(cmd.Args, "-I"+include)
	}

	log.Debugf("compiler: running %s", cmd.Args)
	o, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(err, string(o))
	}

	return nil
}
