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

	fs := afero.NewOsFs()

	rootDir, err := afero.TempDir(fs, "", "btool_builder_test_root")
	if err != nil {
		t.Fatal(err)
	}
	defer fs.RemoveAll(rootDir)

	storeDir, err := afero.TempDir(fs, "", "btool_builder_test_store")
	if err != nil {
		t.Fatal(err)
	}
	defer fs.RemoveAll(storeDir)

	projects := []*testutil.Project{
		testutil.ComplexProjectC(),
		testutil.ComplexProjectCC(),
	}

	for _, project := range projects {
		t.Run(project.Name, func(t *testing.T) {
			project.Root = rootDir
			if err := project.PopulateFS(fs); err != nil {
				t.Fatal(err)
			}

			store := storeDir
			c := compiler.New()
			l := linker.New()
			b := builder.New(fs, project.Root, store, c, l)
			if err := b.Build(project.Graph()); err != nil {
				t.Fatal(err)
			}
		})
	}
}
