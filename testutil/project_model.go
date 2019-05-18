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

var BasicProject = Project{
	Name: "Basic",
	Root: "/tuna/root",
	Nodes: []ProjectNode{
		ProjectNode{
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
		ProjectNode{
			Name: "master.h",
			Includes: []string{
				"<stdlib.h>",
			},
			Dependencies: []string{},
		},
		ProjectNode{
			Name: "dep-0/dep-0a.c",
			Includes: []string{
				"\"dep-0a.h\"",
			},
			Dependencies: []string{
				"dep-0/dep-0a.h",
			},
		},
		ProjectNode{
			Name:         "dep-0/dep-0a.h",
			Includes:     []string{},
			Dependencies: []string{},
		},
	},
}

var BasicProjectWithExtra = Project{
	Name: "BasicWithExtra",
	Root: "/tuna/root",
	Nodes: []ProjectNode{
		ProjectNode{
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
		ProjectNode{
			Name: "master.h",
			Includes: []string{
				"<stdlib.h>",
			},
			Dependencies: []string{},
		},
		ProjectNode{
			Name: "dep-0/dep-0a.c",
			Includes: []string{
				"\"dep-0a.h\"",
			},
			Dependencies: []string{
				"dep-0/dep-0a.h",
			},
		},
		ProjectNode{
			Name:         "dep-0/dep-0a.h",
			Includes:     []string{},
			Dependencies: []string{},
		},
		ProjectNode{
			Name: "dep-1/dep-1a.c",
			Includes: []string{
				"\"dep-0/dep-0a.h\"",
				"\"dep-1a.h\"",
			},
			Dependencies: []string{
				"dep-0/dep-0a.h",
				"dep-1/dep-1a.h",
			},
		},
		ProjectNode{
			Name:         "dep-1/dep-1a.h",
			Includes:     []string{},
			Dependencies: []string{},
		},
	},
}

var ComplexProject = Project{
	Name: "Complex",
	Root: "/tuna/root",
	Nodes: []ProjectNode{
		ProjectNode{
			Name: "main.c",
			Includes: []string{
				"<stdio.h>",
				"\"master.h\"",
				"\"dep-0/dep-0a.h\"",
				"\"dep-1/dep-1a.h\"",
				"\"dep-2/dep-2a.h\"",
			},
			Dependencies: []string{
				"master.h",
				"dep-0/dep-0a.h",
				"dep-1/dep-1a.h",
				"dep-2/dep-2a.h",
			},
			ExtraContent: `
int main(int argc, char *argv[]) {
  return 0;
}
`,
		},

		ProjectNode{
			Name: "master.h",
			Includes: []string{
				"<stdlib.h>",
			},
			Dependencies: []string{},
		},

		ProjectNode{
			Name: "dep-0/dep-0a.c",
			Includes: []string{
				"\"dep-0a.h\"",
			},
			Dependencies: []string{
				"dep-0/dep-0a.h",
			},
		},
		ProjectNode{
			Name:         "dep-0/dep-0a.h",
			Includes:     []string{},
			Dependencies: []string{},
		},

		ProjectNode{
			Name: "dep-1/dep-1a.c",
			Includes: []string{
				"\"dep-0/dep-0a.h\"",
				"\"dep-1a.h\"",
			},
			Dependencies: []string{
				"dep-0/dep-0a.h",
				"dep-1/dep-1a.h",
			},
		},
		ProjectNode{
			Name: "dep-1/dep-1a.h",
			Includes: []string{
				"\"dep-0/dep-0a.h\"",
			},
			Dependencies: []string{
				"dep-0/dep-0a.h",
			},
		},

		ProjectNode{
			Name: "dep-2/dep-2a.c",
			Includes: []string{
				"\"dep-0/dep-0a.h\"",
				"\"dep-2a.h\"",
				"\"dep-2.h\"",
			},
			Dependencies: []string{
				"dep-0/dep-0a.h",
				"dep-2/dep-2a.h",
				"dep-2/dep-2.h",
			},
		},
		ProjectNode{
			Name: "dep-2/dep-2a.h",
			Includes: []string{
				"\"dep-0/dep-0a.h\"",
				"\"dep-2.h\"",
			},
			Dependencies: []string{
				"dep-0/dep-0a.h",
				"dep-2/dep-2.h",
			},
		},
		ProjectNode{
			Name: "dep-2/dep-2b.c",
			Includes: []string{
				"\"dep-0/dep-0a.h\"",
				"\"dep-2b.h\"",
				"\"dep-2.h\"",
			},
			Dependencies: []string{
				"dep-0/dep-0a.h",
				"dep-2/dep-2b.h",
				"dep-2/dep-2.h",
			},
		},
		ProjectNode{
			Name: "dep-2/dep-2b.h",
			Includes: []string{
				"\"dep-0/dep-0a.h\"",
				"\"dep-2.h\"",
			},
			Dependencies: []string{
				"dep-0/dep-0a.h",
				"dep-2/dep-2.h",
			},
		},
		ProjectNode{
			Name: "dep-2/dep-2.h",
			Includes: []string{
				"<stdio.h>",
				"\"dep-0/dep-0a.h\"",
			},
			Dependencies: []string{
				"dep-0/dep-0a.h",
			},
		},
		ProjectNode{
			Name: "dep-2/dep-2-1/dep-2-1.c",
			Includes: []string{
				"\"dep-2/dep-2.h\"",
				"\"dep-2-1.h\"",
			},
			Dependencies: []string{
				"dep-2/dep-2.h",
				"dep-2/dep-2-1/dep-2-1.h",
			},
		},
		ProjectNode{
			Name: "dep-2/dep-2-1/dep-2-1.h",
			Includes: []string{
				"\"dep-2/dep-2.h\"",
			},
			Dependencies: []string{
				"dep-2/dep-2.h",
			},
		},
		ProjectNode{
			Name: "dep-2/dep-2-2/dep-2-2.c",
			Includes: []string{
				"\"dep-2/dep-2.h\"",
				"\"dep-2-2.h\"",
			},
			Dependencies: []string{
				"dep-2/dep-2.h",
				"dep-2/dep-2-2/dep-2-2.h",
			},
		},
		ProjectNode{
			Name: "dep-2/dep-2-2/dep-2-2.h",
			Includes: []string{
				"\"dep-2/dep-2.h\"",
			},
			Dependencies: []string{
				"dep-2/dep-2.h",
			},
		},
	},
}

type ProjectNode struct {
	Name         string
	Includes     []string
	Dependencies []string
	ExtraContent string
}

type Project struct {
	Name  string
	Root  string
	Nodes []ProjectNode
}

func (p *Project) PopulateFS(fs afero.Fs) error {
	for _, node := range p.Nodes {
		content := bytes.NewBuffer([]byte{})
		content.WriteString(fmt.Sprintf("// %s\n", node.Name))
		for _, include := range node.Includes {
			content.WriteString(fmt.Sprintf("\n#include %s", include))
		}
		if node.ExtraContent != "" {
			content.WriteString(node.ExtraContent)
		}

		file := filepath.Join(p.Root, node.Name)
		dir := filepath.Dir(file)
		if err := fs.MkdirAll(dir, 0700); err != nil {
			return errors.Wrap(err, "mkdir "+dir)
		}
		if err := afero.WriteFile(fs, file, content.Bytes(), 0700); err != nil {
			return errors.Wrap(err, "write file "+file)
		}

		logrus.Debugf("created file " + file)
	}

	return nil
}

func (p *Project) Graph() *graph.Graph {
	g := graph.New()

	for _, node := range p.Nodes {
		pathNode := &graph.Node{
			Name: filepath.Join(p.Root, node.Name),
		}
		g.Add(pathNode, nil)

		for _, dependency := range node.Dependencies {
			dependencyNode := &graph.Node{
				Name: filepath.Join(p.Root, dependency),
			}
			g.Add(pathNode, dependencyNode)
		}
	}

	return g
}
