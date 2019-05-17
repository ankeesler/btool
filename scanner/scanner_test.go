package scanner_test

import (
	"bytes"
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/ankeesler/btool/scanner"
	"github.com/ankeesler/btool/scanner/graph"
	"github.com/ankeesler/btool/testutil"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type node struct {
	name         string
	includes     []string
	dependencies []string
}

func TestScan(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(testutil.NewTestingFormatter(t))

	data := []struct {
		name  string
		nodes []node
		paths map[string][]string
	}{
		{
			name: "Basic",
			nodes: []node{
				node{
					name: "/tuna/root/main.c",
					includes: []string{
						"\"master.h\"",
						"\"dep-0/dep-0a.h\"",
					},
					dependencies: []string{
						"/tuna/root/master.h",
						"/tuna/root/dep-0/dep-0a.h",
					},
				},
				node{
					name: "/tuna/root/master.h",
					includes: []string{
						"<stdlib.h>",
					},
					dependencies: []string{},
				},
				node{
					name: "/tuna/root/dep-0/dep-0a.c",
					includes: []string{
						"\"dep-0a.h\"",
					},
					dependencies: []string{
						"/tuna/root/dep-0/dep-0a.h",
					},
				},
				node{
					name:         "/tuna/root/dep-0/dep-0a.h",
					includes:     []string{},
					dependencies: []string{},
				},
			},
		},
		{
			name: "Complex",
			nodes: []node{
				node{
					name: "/tuna/root/main.c",
					includes: []string{
						"<stdio.h>",
						"\"master.h\"",
						"\"dep-0/dep-0a.h\"",
						"\"dep-1/dep-1a.h\"",
						"\"dep-2/dep-2a.h\"",
					},
					dependencies: []string{
						"/tuna/root/master.h",
						"/tuna/root/dep-0/dep-0a.h",
						"/tuna/root/dep-1/dep-1a.h",
						"/tuna/root/dep-2/dep-2a.h",
					},
				},

				node{
					name: "/tuna/root/master.h",
					includes: []string{
						"<stdlib.h>",
					},
					dependencies: []string{},
				},

				node{
					name: "/tuna/root/dep-0/dep-0a.c",
					includes: []string{
						"\"dep-0a.h\"",
					},
					dependencies: []string{
						"/tuna/root/dep-0/dep-0a.h",
					},
				},
				node{
					name:         "/tuna/root/dep-0/dep-0a.h",
					includes:     []string{},
					dependencies: []string{},
				},

				node{
					name: "/tuna/root/dep-1/dep-1a.c",
					includes: []string{
						"\"dep-0/dep-0a.h\"",
						"\"dep-1a.h\"",
					},
					dependencies: []string{
						"/tuna/root/dep-0/dep-0a.h",
						"/tuna/root/dep-1/dep-1a.h",
					},
				},
				node{
					name: "/tuna/root/dep-1/dep-1a.h",
					includes: []string{
						"\"dep-0/dep-0a.h\"",
					},
					dependencies: []string{
						"/tuna/root/dep-0/dep-0a.h",
					},
				},

				node{
					name: "/tuna/root/dep-2/dep-2a.c",
					includes: []string{
						"\"dep-0/dep-0a.h\"",
						"\"dep-2a.h\"",
						"\"dep-2.h\"",
					},
					dependencies: []string{
						"/tuna/root/dep-0/dep-0a.h",
						"/tuna/root/dep-2/dep-2a.h",
						"/tuna/root/dep-2/dep-2.h",
					},
				},
				node{
					name: "/tuna/root/dep-2/dep-2a.h",
					includes: []string{
						"\"dep-0/dep-0a.h\"",
						"\"dep-2.h\"",
					},
					dependencies: []string{
						"/tuna/root/dep-0/dep-0a.h",
						"/tuna/root/dep-2/dep-2.h",
					},
				},
				node{
					name: "/tuna/root/dep-2/dep-2b.c",
					includes: []string{
						"\"dep-0/dep-0a.h\"",
						"\"dep-2b.h\"",
						"\"dep-2.h\"",
					},
					dependencies: []string{
						"/tuna/root/dep-0/dep-0a.h",
						"/tuna/root/dep-2/dep-2b.h",
						"/tuna/root/dep-2/dep-2.h",
					},
				},
				node{
					name: "/tuna/root/dep-2/dep-2b.h",
					includes: []string{
						"\"dep-0/dep-0a.h\"",
						"\"dep-2.h\"",
					},
					dependencies: []string{
						"/tuna/root/dep-0/dep-0a.h",
						"/tuna/root/dep-2/dep-2.h",
					},
				},
				node{
					name: "/tuna/root/dep-2/dep-2.h",
					includes: []string{
						"<stdio.h>",
						"\"dep-0/dep-0a.h\"",
					},
					dependencies: []string{
						"/tuna/root/dep-0/dep-0a.h",
					},
				},
				node{
					name: "/tuna/root/dep-2/dep-2-1/dep-2-1.c",
					includes: []string{
						"\"dep-2/dep-2.h\"",
						"\"dep-2-1.h\"",
					},
					dependencies: []string{
						"/tuna/root/dep-2/dep-2.h",
						"/tuna/root/dep-2/dep-2-1/dep-2-1.h",
					},
				},
				node{
					name: "/tuna/root/dep-2/dep-2-1/dep-2-1.h",
					includes: []string{
						"\"dep-2/dep-2.h\"",
					},
					dependencies: []string{
						"/tuna/root/dep-2/dep-2.h",
					},
				},
				node{
					name: "/tuna/root/dep-2/dep-2-2/dep-2-2.c",
					includes: []string{
						"\"dep-2/dep-2.h\"",
						"\"dep-2-2.h\"",
					},
					dependencies: []string{
						"/tuna/root/dep-2/dep-2.h",
						"/tuna/root/dep-2/dep-2-2/dep-2-2.h",
					},
				},
				node{
					name: "/tuna/root/dep-2/dep-2-2/dep-2-2.h",
					includes: []string{
						"\"dep-2/dep-2.h\"",
					},
					dependencies: []string{
						"/tuna/root/dep-2/dep-2.h",
					},
				},
			},
		},
	}

	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			fs, err := buildFS(datum.nodes)
			if err != nil {
				t.Fatal(err)
			}

			exG := buildGraph(datum.nodes)

			acG, err := scanner.New(fs, "/tuna/root").Scan()
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(exG, acG) {
				t.Fatalf("expected:\nvvv\n%s\n^^^\nactual:\nvvv\n%s\n^^^\n", exG, acG)
			}
		})
	}
}

func buildFS(nodes []node) (afero.Fs, error) {
	fs := afero.NewMemMapFs()

	for _, node := range nodes {
		relName, err := filepath.Rel("/tuna/root", node.name)
		if err != nil {
			return nil, errors.Wrap(err, "rel")
		}

		content := bytes.NewBuffer([]byte{})
		content.WriteString(fmt.Sprintf("// %s\n", relName))
		for _, include := range node.includes {
			content.WriteString(fmt.Sprintf("\n#include %s", include))
		}

		if err := afero.WriteFile(fs, node.name, content.Bytes(), 0600); err != nil {
			return nil, errors.Wrap(err, "write file "+node.name)
		}
	}

	return fs, nil
}

func buildGraph(nodes []node) *graph.Graph {
	g := graph.New()

	for _, node := range nodes {
		pathNode := &graph.Node{
			Name: node.name,
		}
		g.Add(pathNode, nil)

		for _, dependency := range node.dependencies {
			dependencyNode := &graph.Node{
				Name: dependency,
			}
			g.Add(pathNode, dependencyNode)
		}
	}

	return g
}
