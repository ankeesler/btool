package main

import (
	"path/filepath"

	"github.com/ankeesler/btool/testutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

const (
	root = "/Users/ankeesler/developer/btool/fixture" // :(
)

func main() {
	fs := afero.NewOsFs()
	for _, project := range []*testutil.Project{
		testutil.BasicProjectC(),
		testutil.BasicProjectCC(),
		testutil.BasicProjectWithExtraC(),
		testutil.BasicProjectWithExtraCC(),
		testutil.ComplexProjectC(),
		testutil.ComplexProjectCC(),
	} {
		project.Root = filepath.Join(root, project.Name)
		if err := project.PopulateFS(fs); err != nil {
			logrus.Fatal(err)
		}
	}
}
