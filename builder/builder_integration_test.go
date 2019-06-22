package builder_test

import (
	"testing"

	"github.com/ankeesler/btool/builder"
	"github.com/ankeesler/btool/builder/toolchain"
	"github.com/ankeesler/btool/config"
	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/testutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func TestBuildIntegration(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

	fs := afero.NewOsFs()

	rootDir, err := afero.TempDir(fs, "", "btool_builder_test_root")
	if err != nil {
		t.Fatal(err)
	}
	defer fs.RemoveAll(rootDir)

	configDir, err := afero.TempDir(fs, "", "btool_builder_test_config")
	if err != nil {
		t.Fatal(err)
	}
	defer fs.RemoveAll(configDir)

	projects := []*testutil.Project{
		testutil.ComplexProjectC(),
		testutil.ComplexProjectCC(),
		testutil.BigProjectC(),
		testutil.BigProjectCC(),
	}

	for _, project := range projects {
		t.Run(project.Name, func(t *testing.T) {
			project.Root = rootDir
			if err := project.PopulateFS(fs); err != nil {
				t.Fatal(err)
			}

			cache, err := afero.TempDir(fs, "", "builder_integration_test")
			if err != nil {
				t.Fatal(err)
			}
			defer fs.RemoveAll(cache)

			cfg := config.Config{
				Name:  "some-project-name",
				Root:  project.Root,
				Cache: cache,
			}
			tc := toolchain.New("clang", "clang++", "clang")
			b := builder.New(fs, &cfg, tc)
			if err := b.Build(project.Graph()); err != nil {
				t.Fatal(err)
			}
		})
	}
}
