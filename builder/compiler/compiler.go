package compiler

import (
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Compiler struct {
}

func New() *Compiler {
	return &Compiler{}
}

func (c *Compiler) CompileC(output, input, root string) error {
	return c.compile(output, input, root, "clang")
}

func (c *Compiler) CompileCC(output, input, root string) error {
	return c.compile(output, input, root, "clang++")
}

func (c *Compiler) compile(output, input, root, compiler string) error {
	cmd := exec.Command(
		compiler,
		"-c",
		"-o",
		output,
		input,
		"-I"+root,
		"-Wall",
		"-Werror",
		"-O0",
		"-g",
	)

	logrus.Debugf("running compiler command %s", cmd.Args)
	if msg, err := cmd.CombinedOutput(); err != nil {
		return errors.Wrap(err, fmt.Sprintf("run compiler:\n%s", string(msg)))
	}

	return nil
}
