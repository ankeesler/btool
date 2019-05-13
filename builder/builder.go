package builder

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Builder struct {
	store string
}

func New(store string) *Builder {
	return &Builder{store: store}
}

func (b *Builder) Build(target string) error {
	output := filepath.Join(
		b.store,
		filepath.Base(strings.Replace(target, ".c", "", 1)),
	)
	logrus.Debugf("building target '%s' to output '%s'", target, output)

	printed, err := exec.Command(
		"clang",
		"-o",
		output,
		target,
	).CombinedOutput()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("run compiler (%s)", string(printed)))
	}

	return nil
}
