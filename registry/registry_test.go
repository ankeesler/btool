package registry_test

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/registry"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func TestRegistry(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

	fs := afero.NewMemMapFs()

	root := "/some/path/to/root"
	files := []string{"file", "some/path/to/file"}
	for _, file := range files {
		file = filepath.Join(root, file)
		dir := filepath.Dir(file)
		if err := fs.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := afero.WriteFile(fs, file, []byte(file), 0644); err != nil {
			t.Fatal(err)
		}
	}

	r, err := registry.Create(fs, root)
	if err != nil {
		t.Fatal(err)
	}

	data := []struct {
		path     string
		contents string
		exists   bool
	}{
		{
			path: "index.yml",
			contents: `---
files:
- path: file
  sha256: e444140aa1ae3cad9751dbe444c08cb5bf575f1f4bdd80f57eb8724d8d5939f8
- path: some/path/to/file
  sha256: af6bd7ad1d2fff1c2ce0685618696030ca0ea0e534eea98d26e38eb90eb9728b
`,
			exists: true,
		},
		{
			path:     "file",
			contents: "/some/path/to/root/file",
			exists:   true,
		},
		{
			path:     "some/path/to/file",
			contents: "/some/path/to/root/some/path/to/file",
			exists:   true,
		},
		{
			path:   "nope",
			exists: false,
		},
	}
	for _, datum := range data {
		t.Run(datum.path, func(t *testing.T) {
			f := r.Get(datum.path)
			if !datum.exists {
				if f != nil {
					t.Fatalf("expected %s not to exist", datum.path)
				}
				return
			} else if f == nil {
				t.Fatalf("expected %s to exist", datum.path)
			}

			data, err := ioutil.ReadAll(f)
			if err != nil {
				t.Fatal(err)
			}

			if ex, ac := datum.contents, string(data); ex != ac {
				t.Error(ex, "!=", ac)
			}
		})
	}
}
