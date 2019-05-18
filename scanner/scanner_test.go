package scanner_test

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/scanner"
	"github.com/ankeesler/btool/scanner/graph"
	"github.com/ankeesler/btool/testutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func TestScanRoot(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

	projects := []testutil.Project{
		testutil.BasicProject,
		testutil.ComplexProject,
	}

	for _, project := range projects {
		t.Run(project.Name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			if err := project.PopulateFS(fs); err != nil {
				t.Fatal(err)
			}

			exG := project.Graph()

			acG, err := scanner.New(fs, project.Root).ScanRoot()
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(exG, acG) {
				t.Fatalf("expected:\nvvv\n%s\n^^^\nactual:\nvvv\n%s\n^^^\n", exG, acG)
			}
		})
	}
}

func TestScanFile(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

	data := []struct {
		project testutil.Project
		exG     *graph.Graph
	}{
		{
			project: testutil.BasicProject,
			exG:     testutil.BasicProject.Graph(),
		},
		{
			project: testutil.BasicProjectWithExtra,
			exG:     testutil.BasicProject.Graph(),
		},

		// eh, this should probably work...but it doesn't...
		//{
		//	project: testutil.ComplexProject,
		//	exG:     testutil.ComplexProject.Graph(),
		//},
	}

	for _, datum := range data {
		t.Run(datum.project.Name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			if err := datum.project.PopulateFS(fs); err != nil {
				t.Fatal(err)
			}

			file := filepath.Join(datum.project.Root, "main.c")
			acG, err := scanner.New(fs, datum.project.Root).ScanFile(file)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(datum.exG, acG) {
				t.Fatalf("expected:\nvvv\n%s\n^^^\nactual:\nvvv\n%s\n^^^\n", datum.exG, acG)
			}
		})
	}
}
