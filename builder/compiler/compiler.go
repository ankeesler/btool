package compiler

import (
	"bytes"
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

func (c *Compiler) Compile(input io.Reader, output io.Write) error {
	cmd := exec.Command(
		"clang",
		"-c",
		"-o",
		"-",
		"-",
	)

	errBuf := bytes.NewBuffer([]byte{})
	cmd.Stdin = input
	cmd.Stdout = output
	cmd.Stderr = errBuf

	log.Debugf("running compiler command %s", cmd.Args)
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, fmt.Sprintf("run compiler (%s)", errBuf.String()))
	}

	return nil
}
