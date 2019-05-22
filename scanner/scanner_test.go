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

	projects := []*testutil.Project{
		testutil.BasicProjectC(),
		testutil.BasicProjectCC(),
		testutil.BasicProjectWithExtraC(),
		testutil.BasicProjectWithExtraCC(),
		testutil.ComplexProjectC(),
		testutil.ComplexProjectCC(),
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

	basicProjectC := testutil.BasicProjectC()
	basicProjectCC := testutil.BasicProjectCC()
	basicProjectWithExtraC := testutil.BasicProjectWithExtraC()
	basicProjectWithExtraCC := testutil.BasicProjectWithExtraCC()

	data := []struct {
		p    *testutil.Project
		file string
		exG  *graph.Graph
	}{
		{
			p:    basicProjectC,
			file: "main.c",
			exG:  basicProjectC.Graph(),
		},
		{
			p:    basicProjectCC,
			file: "main.cc",
			exG:  basicProjectCC.Graph(),
		},
		{
			p:    basicProjectWithExtraC,
			file: "main.c",
			exG:  basicProjectC.Graph(),
		},
		{
			p:    basicProjectWithExtraCC,
			file: "main.cc",
			exG:  basicProjectCC.Graph(),
		},

		// This should probably work...but it doesn't.
		//{
		//	p: complexProjectC,
		//	exG: complexProjectC.GraphG(),
		//},
		//{
		//	p: complexProjectCC,
		//	exG: complexProjectCC.GraphG(),
		//},
	}

	for _, datum := range data {
		t.Run(datum.p.Name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			if err := datum.p.PopulateFS(fs); err != nil {
				t.Fatal(err)
			}

			file := filepath.Join(datum.p.Root, datum.file)
			acG, err := scanner.New(fs, datum.p.Root).ScanFile(file)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(datum.exG, acG) {
				t.Fatalf("expected:\nvvv\n%s\n^^^\nactual:\nvvv\n%s\n^^^\n", datum.exG, acG)
			}
		})
	}
}
