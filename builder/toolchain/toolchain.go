package toolchain

import (
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Toolchain struct {
	compilerC  string
	compilerCC string
	linker     string
}

func New(compilerC, compilerCC, linker string) *Toolchain {
	return &Toolchain{
		compilerC:  compilerC,
		compilerCC: compilerCC,
		linker:     linker,
	}
}

func (t *Toolchain) CompileC(output, input string, includeDirs []string) error {
	return compile(output, input, includeDirs, t.compilerC)
}

func (t *Toolchain) CompileCC(output, input string, includeDirs []string) error {
	return compile(output, input, includeDirs, t.compilerCC)
}

func (t *Toolchain) Link(output string, inputs []string) error {
	args := append([]string{"-o", output}, inputs...)
	cmd := exec.Command(
		t.linker,
		args...,
	)

	logrus.Debugf("running link command %s", cmd.Args)
	if msg, err := cmd.CombinedOutput(); err != nil {
		return errors.Wrap(err, fmt.Sprintf("run linker\n%s", string(msg)))
	}

	return nil
}

func compile(
	output, input string,
	includeDirs []string,
	compiler string,
) error {
	cmd := exec.Command(
		compiler,
		"-c",
		"-o",
		output,
		input,
		"-Wall",
		"-Werror",
		"-O0",
		"-g",
	)
	for _, includeDir := range includeDirs {
		cmd.Args = append(cmd.Args, fmt.Sprintf("-I%s", includeDir))
	}

	logrus.Debugf("running compiler command %s", cmd.Args)
	if msg, err := cmd.CombinedOutput(); err != nil {
		return errors.Wrap(err, fmt.Sprintf("run compiler:\n%s", string(msg)))
	}

	return nil
}
