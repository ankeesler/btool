package builder

import (
	"path/filepath"

	"github.com/pkg/errors"
)

func (b *Builder) Clean() error {
	if err := b.fs.RemoveAll(
		filepath.Join(
			b.store,
			"objects",
		),
	); err != nil {
		return errors.Wrap(err, "remove objects")
	}

	if err := b.fs.RemoveAll(
		filepath.Join(
			b.store,
			"binaries",
		),
	); err != nil {
		return errors.Wrap(err, "remove binaries")
	}

	return nil
}
