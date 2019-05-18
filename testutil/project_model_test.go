package testutil_test

import (
	"testing"

	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/testutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func TestComplexProject(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

	for _, project := range []testutil.Project{
		testutil.BasicProject,
		testutil.BasicProjectWithExtra,
		testutil.ComplexProject,
	} {
		t.Run(project.Name, func(t *testing.T) {
			if g := testutil.ComplexProject.Graph(); g == nil {
				t.Errorf("expected g to not be nil")
			}

			if err := testutil.ComplexProject.PopulateFS(afero.NewMemMapFs()); err != nil {
				t.Error(err)
			}
		})
	}
}
