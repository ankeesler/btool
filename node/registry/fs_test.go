package registry_test

import (
	"path/filepath"
	"testing"

	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/node/registry"
	"github.com/ankeesler/btool/node/registry/testutil"
	"github.com/go-test/deep"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

func TestFSRegistry(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

	fs := afero.NewMemMapFs()

	root := "/some/path/to/root"
	files := map[string][]*registry.Node{
		"file_a_btool.yml":              testutil.FileANodes(),
		"some/path/to/file_b_btool.yml": testutil.FileBNodes(),
		"whatever.yml":                  nil,
	}
	for file, nodes := range files {
		file = filepath.Join(root, file)
		dir := filepath.Dir(file)
		if err := fs.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}

		data, err := yaml.Marshal(nodes)
		if err != nil {
			t.Fatal(err)
		}

		if err := afero.WriteFile(fs, file, data, 0644); err != nil {
			t.Fatal(err)
		}
	}

	r, err := registry.CreateFSRegistry(fs, root)
	if err != nil {
		t.Fatal(err)
	}

	exIndex := testutil.Index()
	acIndex, err := r.Index()
	if err != nil {
		t.Error(err)
	} else if diff := deep.Equal(exIndex, acIndex); diff != nil {
		t.Error(diff)
	}

	data := []struct {
		name   string
		nodes  []*registry.Node
		exists bool
	}{
		{
			name:   "file_a_btool.yml",
			nodes:  testutil.FileANodes(),
			exists: true,
		},
		{
			name:   "some/path/to/file_b_btool.yml",
			nodes:  testutil.FileBNodes(),
			exists: true,
		},
		{
			name:   "whatever.yml",
			exists: false,
		},
		{
			name:   "nope",
			exists: false,
		},
	}
	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			acNodes, err := r.Nodes(datum.name)
			if err != nil {
				t.Fatal(err)
			}

			if datum.exists {
				if acNodes == nil {
					t.Fatal("expected nodes to exist")
				}
			} else {
				if acNodes != nil {
					t.Fatal("expected nodes not to exist")
				}
				return
			}

			if diff := deep.Equal(datum.nodes, acNodes); diff != nil {
				t.Error(diff)
			}
		})
	}
}
