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

func (c *Compiler) Compile(output, input string) error {
	cmd := exec.Command(
		"clang",
		"-c",
		"-o",
		output,
		input,
	)

	logrus.Debugf("running compiler command %s", cmd.Args)
	if stderr, err := cmd.CombinedOutput(); err != nil {
		return errors.Wrap(err, fmt.Sprintf("run compiler:\n%s", string(stderr)))
	}

	return nil
}
