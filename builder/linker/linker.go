package linker

import (
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Linker struct {
}

func New() *Linker {
	return &Linker{}
}

func (c *Linker) Link(output string, inputs []string) error {
	args := append([]string{"-o", output}, inputs...)
	cmd := exec.Command(
		"clang",
		args...,
	)

	logrus.Debugf("running link command %s", cmd.Args)
	if msg, err := cmd.CombinedOutput(); err != nil {
		return errors.Wrap(err, fmt.Sprintf("run linker\n%s", string(msg)))
	}

	return nil
}
