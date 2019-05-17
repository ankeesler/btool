package linker

import (
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Linker struct {
}

func New() *Linker {
	return &Linker{}
}

func (c *Linker) Link(output string, input []string) error {
	args := append([]string{"-o", output}, input...)
	cmd := exec.Command(
		"clang",
		args...,
	)
	log.Debugf("running linker command %s", cmd.Args)

	if stdoutAndErr, err := cmd.CombinedOutput(); err != nil {
		return errors.Wrap(err, fmt.Sprintf("run linker (%s)", string(stdoutAndErr)))
	}

	return nil
}
