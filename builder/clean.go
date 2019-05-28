package builder

import (
	"github.com/pkg/errors"
)

func (b *Builder) Clean() error {
	if err := b.fs.RemoveAll(b.objectsDir()); err != nil {
		return errors.Wrap(err, "remove objects")
	}

	if err := b.fs.RemoveAll(b.binariesDir()); err != nil {
		return errors.Wrap(err, "remove binaries")
	}

	return nil
}
