package linker

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Linker struct {
}

func New() *Linker {
	return &Linker{}
}

func (c *Linker) Link(output io.Writer, inputs []io.Reader) error {
	errBuf := bytes.NewBuffer([]byte{})

	cmd := exec.Command(
		"clang",
		"-c",
		"-o",
		"-",
		"-",
	)
	cmd.Stdin = io.MultiReader(inputs...)
	cmd.Stdout = output
	cmd.Stderr = errBuf

	logrus.Debugf("running link command %s", cmd.Args)
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, fmt.Sprintf("run linker (%s)", errBuf.String()))
	}

	return nil
}
