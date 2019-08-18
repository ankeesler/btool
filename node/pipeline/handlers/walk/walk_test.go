package walk_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ankeesler/btool/node/pipeline/handlers/walk"
	"github.com/go-test/deep"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

// Note! afero does not contain support for symlinking currently.
// Therefore we will run this test with an OS FS (as opposed to a memory FS).
// See https://github.com/spf13/afero/pull/212/files.

func TestWalk(t *testing.T) {
	fs := afero.NewOsFs()
	files := []string{
		filepath.Join("dir-a0"),
		filepath.Join("dir-a0", "file-a0a.txt"),
		filepath.Join("dir-a0", "dir-a1"),
		filepath.Join("dir-a0", "dir-a1", "file-a1a.tuna"),
		filepath.Join("dir-a0", "dir-a1", "dir-a2"),

		filepath.Join("dir-b0"),

		filepath.Join("dir-c0"),
		filepath.Join("dir-c0", "dir-c1"),
		filepath.Join("dir-c0", "dir-c1", "file-c2a.txt"),
		filepath.Join("dir-c0", "dir-c1", "file-c3b.fish"),

		filepath.Join("dir-d0"),
		filepath.Join("dir-d0", "symlink.dir-c0_dir-c1"),

		filepath.Join("dir-e0"),
		filepath.Join("dir-e0", "symlink.dir-e0"),
	}
	root, err := ioutil.TempDir("", "btool_walk_test")
	require.Nil(t, err)
	defer func() {
		require.Nil(t, fs.RemoveAll(root))
	}()
	populateFS(t, fs, files, root)

	data := []struct {
		name    string
		root    string
		exts    []string
		failure bool
		paths   []string
	}{
		{
			name:    "Basic",
			root:    "",
			exts:    []string{".txt", ".fish"},
			failure: false,
			paths: []string{
				filepath.Join("dir-a0", "file-a0a.txt"),
				filepath.Join("dir-c0", "dir-c1", "file-c2a.txt"),
				filepath.Join("dir-c0", "dir-c1", "file-c3b.fish"),
			},
		},
		{
			name:    "Failure",
			root:    "",
			exts:    []string{".txt", ".fish"},
			failure: true,
			paths: []string{
				filepath.Join("dir-a0", "file-a0a.txt"),
			},
		},
		{
			name:    "Symlink",
			root:    "dir-d0",
			exts:    []string{".fish"},
			failure: false,
			paths: []string{
				filepath.Join("dir-d0", "symlink.dir-c0_dir-c1", "file-c3b.fish"),
			},
		},
	}

	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			paths := make([]string, 0)
			walkFn := func(path string) error {
				relPath, err := filepath.Rel(root, path)
				require.Nil(t, err)
				paths = append(paths, relPath)

				if datum.failure {
					return errors.New("failure")
				} else {
					return nil
				}
			}

			err := walk.Walk(fs, filepath.Join(root, datum.root), datum.exts, walkFn)
			if datum.failure {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
				require.Nil(t, deep.Equal(datum.paths, paths))
			}
		})
	}
}

func populateFS(t *testing.T, fs afero.Fs, files []string, root string) {
	for _, file := range files {
		base := filepath.Base(file)
		if strings.HasPrefix(base, "dir") {
			require.Nil(
				t,
				fs.Mkdir(
					filepath.Join(root, file),
					0755,
				),
			)
		} else if strings.HasPrefix(base, "file") {
			require.Nil(
				t,
				afero.WriteFile(
					fs,
					filepath.Join(root, file),
					[]byte(file),
					0644,
				),
			)
		} else if strings.HasPrefix(base, "symlink") {
			oldFile := strings.NewReplacer(
				"symlink.", "",
				"_", "/",
			).Replace(base)
			require.Nil(
				t,
				os.Symlink(
					filepath.Join(root, oldFile),
					filepath.Join(root, file),
				),
			)
		}
	}
}