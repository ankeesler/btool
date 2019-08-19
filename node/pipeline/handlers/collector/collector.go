// Package collector provides filesystem walking functionality that is specific
// to btool.
package collector

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ankeesler/btool/log"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// Collector is an object that can walk a file tree and collect file names.
type Collector struct {
}

// New creates a new Collector.
func New() *Collector {
	return &Collector{}
}

// Note! afero does not contain support for symlinking currently.
// See https://github.com/spf13/afero/pull/212/files.

// Collect will walk a filesystem from the provided root. It has opinionated
// specifics (for btool):
//   1. It only passes files to the provided walkFn
//   2. It follows symlinks...without much loop detection :(
//   3. It quits immediately upon error
//   4. It only passes files with the provided file extensions to the provided
//      walkFn
//   5. The file paths passed to the walkFn contain the symlink in them
func (c *Collector) Collect(
	root string,
	exts []string,
) ([]string, error) {
	paths := make([]string, 0)
	walkFn := func(path string) error {
		paths = append(paths, path)
		return nil
	}

	fs := afero.NewOsFs()
	visited := make(map[string]bool)
	if err := walk(fs, root, "", exts, walkFn, visited); err != nil {
		return nil, errors.Wrap(err, "walk")
	}

	return paths, nil
}

func walk(
	fs afero.Fs,
	root string,
	linkRoot string,
	exts []string,
	walkFn func(string) error,
	visited map[string]bool,
) error {
	log.Debugf("walk root %s (link root: %s) for exts %s", root, linkRoot, exts)
	return afero.Walk(
		fs,
		root,
		func(path string, info os.FileInfo, err error) error {
			if visited[path] {
				return nil
			}
			visited[path] = true

			if err != nil {
				return err
			}

			if info.IsDir() {
				log.Debugf("skipping directory %s", path)
			} else if (info.Mode() & os.ModeSymlink) != 0 {
				log.Debugf("looking at link %s", path)
				link, err := os.Readlink(path)
				if err != nil {
					return errors.Wrap(err, "read link")
				}

				return walk(fs, link, path, exts, walkFn, visited)
			} else {
				log.Debugf("looking at file %s", path)

				var realPath string
				if linkRoot != "" {
					realPath = strings.ReplaceAll(path, root, linkRoot)
					log.Debugf("actually though %s", realPath)
				} else {
					realPath = path
				}

				actualExt := filepath.Ext(path)
				for _, ext := range exts {
					if ext == actualExt {
						err := walkFn(realPath)
						return err
					}
				}
			}

			return nil
		},
	)
}
