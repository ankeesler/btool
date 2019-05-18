package builder_test

import (
	"testing"

	"github.com/ankeesler/btool/builder"
	"github.com/ankeesler/btool/builder/compiler"
	"github.com/ankeesler/btool/builder/linker"
	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/testutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func TestBuild(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

	fs := afero.NewMemMapFs()
	if err := testutil.ComplexProject.PopulateFS(fs); err != nil {
		t.Fatal(err)
	}

	root := "/tuna/root"
	store := "/tmp/btool-store"
	c := compiler.New()
	l := linker.New()
	b := builder.New(fs, root, store, c, l)
	if err := b.Build(testutil.ComplexProject.Graph()); err != nil {
		t.Fatal(err)
	}
}
