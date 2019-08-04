package resolvers_test

import (
	"testing"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/resolvers"
	"github.com/spf13/afero"
)

func TestGenerateC(t *testing.T) {
	fs := afero.NewMemMapFs()
	dir := "some/dir"
	gc := resolvers.NewGenerateC(fs, dir)

	n := node.New("some/path/to/file.c")
	if err := gc.Resolve(n); err != nil {
		t.Fatal(err)
	}

	data, err := afero.ReadFile(fs, "some/dir/some/path/to/file.c")
	if err != nil {
		t.Fatal(err)
	}

	if ex, ac := "#include \"file.h\"\n", string(data); ex != ac {
		t.Error(ex, "!=", ac)
	}
}
