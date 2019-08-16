// Package walk provides filesystem walking functionality that is specific to
// btool.
package walk

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// Note! afero does not contain support for symlinking currently.
// See https://github.com/spf13/afero/pull/212/files.

// Walk will walk a filesystem from the provided root. It has opinionated
// specifics (for btool):
//   1. It only passes files to the provided walkFn
//   2. It follows symlinks...without loop detection :(
//   3. It quits immediately upon error
//   4. It only passes files with the provided file extensions to the provided
//      walkFn
//   5. The file paths passed to the walkFn contain the symlink in them
func Walk(
	fs afero.Fs,
	root string,
	exts []string,
	walkFn func(string) error,
) error {
	return walk(fs, root, "", exts, walkFn)
}

func walk(
	fs afero.Fs,
	root string,
	linkRoot string,
	exts []string,
	walkFn func(string) error,
) error {
	logrus.Debugf("walk root %s (link root %s) for exts %s", root, linkRoot, exts)
	return afero.Walk(
		fs,
		root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				logrus.Debugf("skipping directory %s", path)
			} else if (info.Mode() & os.ModeSymlink) != 0 {
				logrus.Debugf("looking at link %s", path)
				link, err := os.Readlink(path)
				if err != nil {
					return errors.Wrap(err, "read link")
				}

				return walk(fs, link, path, exts, walkFn)
			} else {
				logrus.Debugf("looking at file %s", path)

				if linkRoot != "" {
					path = strings.ReplaceAll(path, root, linkRoot)
					logrus.Debugf("actually though %s", path)
				}

				actualExt := filepath.Ext(path)
				for _, ext := range exts {
					if ext == actualExt {
						return walkFn(path)
					}
				}
			}

			return nil
		},
	)
}
