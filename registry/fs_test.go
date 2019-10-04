package registry_test

import (
	"path/filepath"
	"testing"

	"github.com/ankeesler/btool/registry"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestFSRegistry(t *testing.T) {
	fs := afero.NewMemMapFs()

	root := "/some/path/to/root"
	files := map[string]*registry.Gaggle{
		"file_a_btool.yml":              fileAGaggle(),
		"some/path/to/file_b_btool.yml": fileBGaggle(),
	}
	for file, gaggle := range files {
		file = filepath.Join(root, file)
		dir := filepath.Dir(file)
		if err := fs.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}

		data, err := yaml.Marshal(gaggle)
		if err != nil {
			t.Fatal(err)
		}

		if err := afero.WriteFile(fs, file, data, 0644); err != nil {
			t.Fatal(err)
		}
	}

	r, err := registry.CreateFSRegistry(fs, root, "some-index")
	if err != nil {
		t.Fatal(err)
	}

	exIndex := index()
	acIndex, err := r.Index()
	require.Nil(t, err)
	require.Equal(t, exIndex, acIndex)

	data := []struct {
		name   string
		gaggle *registry.Gaggle
		exists bool
	}{
		{
			name:   "file_a_btool.yml",
			gaggle: fileAGaggle(),
			exists: true,
		},
		{
			name:   "some/path/to/file_b_btool.yml",
			gaggle: fileBGaggle(),
			exists: true,
		},
		{
			name:   "index.yml",
			exists: false,
		},
		{
			name:   "nope",
			exists: false,
		},
	}
	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			acGaggle, err := r.Gaggle(datum.name)
			if err != nil {
				t.Fatal(err)
			}

			if datum.exists {
				if acGaggle == nil {
					t.Fatal("expected gaggle to exist")
				}
			} else {
				if acGaggle != nil {
					t.Fatal("expected gaggle not to exist")
				}
				return
			}

			assert.Equal(t, datum.gaggle, acGaggle)
		})
	}
}
