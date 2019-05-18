package testutil

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/ankeesler/btool/scanner/graph"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

var ComplexProject = Project{
	Name: "Complex",
	Root: "/tuna/root",
	Nodes: []ProjectNode{
		ProjectNode{
			Name: "/tuna/root/main.c",
			Includes: []string{
				"<stdio.h>",
				"\"master.h\"",
				"\"dep-0/dep-0a.h\"",
				"\"dep-1/dep-1a.h\"",
				"\"dep-2/dep-2a.h\"",
			},
			Dependencies: []string{
				"/tuna/root/master.h",
				"/tuna/root/dep-0/dep-0a.h",
				"/tuna/root/dep-1/dep-1a.h",
				"/tuna/root/dep-2/dep-2a.h",
			},
		},

		ProjectNode{
			Name: "/tuna/root/master.h",
			Includes: []string{
				"<stdlib.h>",
			},
			Dependencies: []string{},
		},

		ProjectNode{
			Name: "/tuna/root/dep-0/dep-0a.c",
			Includes: []string{
				"\"dep-0a.h\"",
			},
			Dependencies: []string{
				"/tuna/root/dep-0/dep-0a.h",
			},
		},
		ProjectNode{
			Name:         "/tuna/root/dep-0/dep-0a.h",
			Includes:     []string{},
			Dependencies: []string{},
		},

		ProjectNode{
			Name: "/tuna/root/dep-1/dep-1a.c",
			Includes: []string{
				"\"dep-0/dep-0a.h\"",
				"\"dep-1a.h\"",
			},
			Dependencies: []string{
				"/tuna/root/dep-0/dep-0a.h",
				"/tuna/root/dep-1/dep-1a.h",
			},
		},
		ProjectNode{
			Name: "/tuna/root/dep-1/dep-1a.h",
			Includes: []string{
				"\"dep-0/dep-0a.h\"",
			},
			Dependencies: []string{
				"/tuna/root/dep-0/dep-0a.h",
			},
		},

		ProjectNode{
			Name: "/tuna/root/dep-2/dep-2a.c",
			Includes: []string{
				"\"dep-0/dep-0a.h\"",
				"\"dep-2a.h\"",
				"\"dep-2.h\"",
			},
			Dependencies: []string{
				"/tuna/root/dep-0/dep-0a.h",
				"/tuna/root/dep-2/dep-2a.h",
				"/tuna/root/dep-2/dep-2.h",
			},
		},
		ProjectNode{
			Name: "/tuna/root/dep-2/dep-2a.h",
			Includes: []string{
				"\"dep-0/dep-0a.h\"",
				"\"dep-2.h\"",
			},
			Dependencies: []string{
				"/tuna/root/dep-0/dep-0a.h",
				"/tuna/root/dep-2/dep-2.h",
			},
		},
		ProjectNode{
			Name: "/tuna/root/dep-2/dep-2b.c",
			Includes: []string{
				"\"dep-0/dep-0a.h\"",
				"\"dep-2b.h\"",
				"\"dep-2.h\"",
			},
			Dependencies: []string{
				"/tuna/root/dep-0/dep-0a.h",
				"/tuna/root/dep-2/dep-2b.h",
				"/tuna/root/dep-2/dep-2.h",
			},
		},
		ProjectNode{
			Name: "/tuna/root/dep-2/dep-2b.h",
			Includes: []string{
				"\"dep-0/dep-0a.h\"",
				"\"dep-2.h\"",
			},
			Dependencies: []string{
				"/tuna/root/dep-0/dep-0a.h",
				"/tuna/root/dep-2/dep-2.h",
			},
		},
		ProjectNode{
			Name: "/tuna/root/dep-2/dep-2.h",
			Includes: []string{
				"<stdio.h>",
				"\"dep-0/dep-0a.h\"",
			},
			Dependencies: []string{
				"/tuna/root/dep-0/dep-0a.h",
			},
		},
		ProjectNode{
			Name: "/tuna/root/dep-2/dep-2-1/dep-2-1.c",
			Includes: []string{
				"\"dep-2/dep-2.h\"",
				"\"dep-2-1.h\"",
			},
			Dependencies: []string{
				"/tuna/root/dep-2/dep-2.h",
				"/tuna/root/dep-2/dep-2-1/dep-2-1.h",
			},
		},
		ProjectNode{
			Name: "/tuna/root/dep-2/dep-2-1/dep-2-1.h",
			Includes: []string{
				"\"dep-2/dep-2.h\"",
			},
			Dependencies: []string{
				"/tuna/root/dep-2/dep-2.h",
			},
		},
		ProjectNode{
			Name: "/tuna/root/dep-2/dep-2-2/dep-2-2.c",
			Includes: []string{
				"\"dep-2/dep-2.h\"",
				"\"dep-2-2.h\"",
			},
			Dependencies: []string{
				"/tuna/root/dep-2/dep-2.h",
				"/tuna/root/dep-2/dep-2-2/dep-2-2.h",
			},
		},
		ProjectNode{
			Name: "/tuna/root/dep-2/dep-2-2/dep-2-2.h",
			Includes: []string{
				"\"dep-2/dep-2.h\"",
			},
			Dependencies: []string{
				"/tuna/root/dep-2/dep-2.h",
			},
		},
	},
}

type ProjectNode struct {
	Name         string
	Includes     []string
	Dependencies []string
}

type Project struct {
	Name  string
	Root  string
	Nodes []ProjectNode
}

func (p *Project) PopulateFS(fs afero.Fs) error {
	for _, node := range p.Nodes {
		relName, err := filepath.Rel(p.Root, node.Name)
		if err != nil {
			return errors.Wrap(err, "rel")
		}

		content := bytes.NewBuffer([]byte{})
		content.WriteString(fmt.Sprintf("// %s\n", relName))
		for _, include := range node.Includes {
			content.WriteString(fmt.Sprintf("\n#include %s", include))
		}

		if err := afero.WriteFile(fs, node.Name, content.Bytes(), 0600); err != nil {
			return errors.Wrap(err, "write file "+node.Name)
		}

		logrus.Debugf("created file " + node.Name)
	}

	return nil
}

func (p *Project) Graph() *graph.Graph {
	g := graph.New()

	for _, node := range p.Nodes {
		pathNode := &graph.Node{
			Name: node.Name,
		}
		g.Add(pathNode, nil)

		for _, dependency := range node.Dependencies {
			dependencyNode := &graph.Node{
				Name: dependency,
			}
			g.Add(pathNode, dependencyNode)
		}
	}

	return g
}
