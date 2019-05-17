package compiler

import (
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Compiler struct {
}

func New() *Compiler {
	return &Compiler{}
}

func (c *Compiler) Compile(input, output string) error {
	cmd := exec.Command(
		"clang",
		"-c",
		"-o",
		output,
		input,
	)
	log.Debugf("running compiler command %s", cmd.Args)

	if stdoutAndErr, err := cmd.CombinedOutput(); err != nil {
		return errors.Wrap(err, fmt.Sprintf("run compiler (%s)", string(stdoutAndErr)))
	}

	return nil
}
