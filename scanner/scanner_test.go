package scanner_test

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/ankeesler/btool/config"
	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/scanner"
	"github.com/ankeesler/btool/scanner/graph"
	"github.com/ankeesler/btool/scanner/scannerfakes"
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
		testutil.BigProjectC(),
		testutil.BigProjectCC(),
	}

	for _, project := range projects {
		t.Run(project.Name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			if err := project.PopulateFS(fs); err != nil {
				t.Fatal(err)
			}

			exG := project.Graph()

			c := config.Config{
				Name:  "some-project-name",
				Root:  project.Root,
				Cache: "/some/cache/root",
			}
			r := &scannerfakes.FakeResolver{}
			acG, err := scanner.New(fs, &c, r).ScanRoot()
			if err != nil {
				t.Fatal(err)
			}

			if err := graph.Equal(exG, acG); err != nil {
				t.Fatal(err)
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
	dependencyProject := &testutil.Project{
		Name: "Dependency",
		Root: "/tuna/root",
		Nodes: []*testutil.ProjectNode{
			&testutil.ProjectNode{
				Name: "main.c",
				Includes: []string{
					"<stdlib.h>",
					"\"dep-0/dep-0.h\"",
				},
				Dependencies: []string{
					"dep-0/dep-0.h",
				},
			},
			&testutil.ProjectNode{
				Name:         "dep-0/dep-0.h",
				Includes:     []string{},
				Dependencies: []string{},
			},
			&testutil.ProjectNode{
				Name: "dep-0/dep-0.c",
				Includes: []string{
					"\"dep-0.h\"",
					"\"some/path/to/dep.h\"",
				},
				Dependencies: []string{
					"dep-0/dep-0.h",
					"/cache/dependencies/dep/include/dep.h",
				},
			},
		},
	}

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
		{
			p:    dependencyProject,
			file: "main.c",
			exG:  dependencyProject.Graph(),
		},

		// This should probably work...but it doesn't.
		//{
		//	p: complexProjectC,
		//	exG: complexProjectC.Graph(),
		//},
		//{
		//	p: complexProjectCC,
		//	exG: complexProjectCC.Graph(),
		//},
		//{
		//	p: bigProjectC,
		//	exG: bigProjectC.Graph(),
		//},
		//{
		//	p: bigProjectCC,
		//	exG: bigProjectCC.Graph(),
		//},
	}

	for _, datum := range data {
		t.Run(datum.p.Name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			if err := datum.p.PopulateFS(fs); err != nil {
				t.Fatal(err)
			}

			r := &scannerfakes.FakeResolver{}
			r.ResolveIncludeStub = func(include string) (string, error) {
				logrus.Debug("ResolveInclude:", include)
				if include == "some/path/to/dep.h" {
					return "/cache/dependencies/dep/include/dep.h", nil
				} else {
					panic("should not get here!")
				}
			}
			r.ResolveSourcesStub = func(include string) ([]string, error) {
				logrus.Debug("ResolveSources:", include)
				if include == "some/path/to/dep.h" {
					return []string{
						"/cache/dependencies/dep/source/dep-a.c",
						"/cache/dependencies/dep/source/dep-b.c",
					}, nil
				} else {
					return []string{}, nil
				}
			}

			c := config.Config{
				Name:  "some-project-name",
				Root:  datum.p.Root,
				Cache: "/some/cache/root",
			}
			file := filepath.Join(datum.p.Root, datum.file)
			acG, err := scanner.New(fs, &c, r).ScanFile(file)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(datum.exG, acG) {
				t.Fatalf("expected:\nvvv\n%s\n^^^\nactual:\nvvv\n%s\n^^^\n", datum.exG, acG)
			}
		})
	}
}
