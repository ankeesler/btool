package scanner_test

import (
	"reflect"
	"testing"

	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/scanner"
	"github.com/ankeesler/btool/testutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func TestScan(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

	projects := []testutil.Project{
		testutil.Project{
			Name: "Basic",
			Root: "/tuna/root",
			Nodes: []testutil.ProjectNode{
				testutil.ProjectNode{
					Name: "main.c",
					Includes: []string{
						"\"master.h\"",
						"\"dep-0/dep-0a.h\"",
					},
					Dependencies: []string{
						"master.h",
						"dep-0/dep-0a.h",
					},
				},
				testutil.ProjectNode{
					Name: "master.h",
					Includes: []string{
						"<stdlib.h>",
					},
					Dependencies: []string{},
				},
				testutil.ProjectNode{
					Name: "dep-0/dep-0a.c",
					Includes: []string{
						"\"dep-0a.h\"",
					},
					Dependencies: []string{
						"dep-0/dep-0a.h",
					},
				},
				testutil.ProjectNode{
					Name:         "dep-0/dep-0a.h",
					Includes:     []string{},
					Dependencies: []string{},
				},
			},
		},
		testutil.ComplexProject,
	}

	for _, project := range projects {
		t.Run(project.Name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			if err := project.PopulateFS(fs); err != nil {
				t.Fatal(err)
			}

			exG := project.Graph()

			acG, err := scanner.New(fs, project.Root).Scan()
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(exG, acG) {
				t.Fatalf("expected:\nvvv\n%s\n^^^\nactual:\nvvv\n%s\n^^^\n", exG, acG)
			}
		})
	}
}
