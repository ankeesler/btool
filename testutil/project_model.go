package testutil

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ankeesler/btool/scanner/graph"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

var basicProject = Project{
	Name: "Basic",
	Root: "/tuna/root",
	Nodes: []*ProjectNode{
		&ProjectNode{
			Name: "main.FILE_EXTENSION",
			Includes: []string{
				"\"master.h\"",
				"\"dep-0/dep-0a.h\"",
			},
			Dependencies: []string{
				"master.h",
				"dep-0/dep-0a.h",
			},
		},
		&ProjectNode{
			Name: "master.h",
			Includes: []string{
				"<stdlib.h>",
			},
			Dependencies: []string{},
		},
		&ProjectNode{
			Name: "dep-0/dep-0a.FILE_EXTENSION",
			Includes: []string{
				"\"dep-0a.h\"",
			},
			Dependencies: []string{
				"dep-0/dep-0a.h",
			},
		},
		&ProjectNode{
			Name:         "dep-0/dep-0a.h",
			Includes:     []string{},
			Dependencies: []string{},
		},
	},
}

var basicProjectWithExtra = Project{
	Name: "BasicWithExtra",
	Root: "/tuna/root",
	Nodes: []*ProjectNode{
		&ProjectNode{
			Name: "main.FILE_EXTENSION",
			Includes: []string{
				"\"master.h\"",
				"\"dep-0/dep-0a.h\"",
			},
			Dependencies: []string{
				"master.h",
				"dep-0/dep-0a.h",
			},
		},
		&ProjectNode{
			Name: "master.h",
			Includes: []string{
				"<stdlib.h>",
			},
			Dependencies: []string{},
		},
		&ProjectNode{
			Name: "dep-0/dep-0a.FILE_EXTENSION",
			Includes: []string{
				"\"dep-0a.h\"",
			},
			Dependencies: []string{
				"dep-0/dep-0a.h",
			},
		},
		&ProjectNode{
			Name:         "dep-0/dep-0a.h",
			Includes:     []string{},
			Dependencies: []string{},
		},
		&ProjectNode{
			Name: "dep-1/dep-1a.FILE_EXTENSION",
			Includes: []string{
				"\"dep-0/dep-0a.h\"",
				"\"dep-1a.h\"",
			},
			Dependencies: []string{
				"dep-0/dep-0a.h",
				"dep-1/dep-1a.h",
			},
		},
		&ProjectNode{
			Name:         "dep-1/dep-1a.h",
			Includes:     []string{},
			Dependencies: []string{},
		},
	},
}

var complexProject = Project{
	Name: "Complex",
	Root: "/tuna/root",
	Nodes: []*ProjectNode{
		&ProjectNode{
			Name: "main.FILE_EXTENSION",
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
  printf("hey! i am running!\n");
  return 0;
}
`,
		},

		&ProjectNode{
			Name: "master.h",
			Includes: []string{
				"<stdlib.h>",
			},
			Dependencies: []string{},
		},

		&ProjectNode{
			Name: "dep-0/dep-0a.FILE_EXTENSION",
			Includes: []string{
				"\"dep-0a.h\"",
			},
			Dependencies: []string{
				"dep-0/dep-0a.h",
			},
		},
		&ProjectNode{
			Name:         "dep-0/dep-0a.h",
			Includes:     []string{},
			Dependencies: []string{},
		},

		&ProjectNode{
			Name: "dep-1/dep-1a.FILE_EXTENSION",
			Includes: []string{
				"\"dep-0/dep-0a.h\"",
				"\"dep-1a.h\"",
			},
			Dependencies: []string{
				"dep-0/dep-0a.h",
				"dep-1/dep-1a.h",
			},
		},
		&ProjectNode{
			Name: "dep-1/dep-1a.h",
			Includes: []string{
				"\"dep-0/dep-0a.h\"",
			},
			Dependencies: []string{
				"dep-0/dep-0a.h",
			},
		},

		&ProjectNode{
			Name: "dep-2/dep-2a.FILE_EXTENSION",
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
		&ProjectNode{
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
		&ProjectNode{
			Name: "dep-2/dep-2b.FILE_EXTENSION",
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
		&ProjectNode{
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
		&ProjectNode{
			Name: "dep-2/dep-2.h",
			Includes: []string{
				"<stdio.h>",
				"\"dep-0/dep-0a.h\"",
			},
			Dependencies: []string{
				"dep-0/dep-0a.h",
			},
		},
		&ProjectNode{
			Name: "dep-2/dep-2-1/dep-2-1.FILE_EXTENSION",
			Includes: []string{
				"\"dep-2/dep-2.h\"",
				"\"dep-2-1.h\"",
			},
			Dependencies: []string{
				"dep-2/dep-2.h",
				"dep-2/dep-2-1/dep-2-1.h",
			},
		},
		&ProjectNode{
			Name: "dep-2/dep-2-1/dep-2-1.h",
			Includes: []string{
				"\"dep-2/dep-2.h\"",
			},
			Dependencies: []string{
				"dep-2/dep-2.h",
			},
		},
		&ProjectNode{
			Name: "dep-2/dep-2-2/dep-2-2.FILE_EXTENSION",
			Includes: []string{
				"\"dep-2/dep-2.h\"",
				"\"dep-2-2.h\"",
			},
			Dependencies: []string{
				"dep-2/dep-2.h",
				"dep-2/dep-2-2/dep-2-2.h",
			},
		},
		&ProjectNode{
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

func BasicProjectC() *Project {
	return executeTemplate(deepCopy(&basicProject), "c")
}

func BasicProjectCC() *Project {
	return executeTemplate(deepCopy(&basicProject), "cc")
}

func BasicProjectWithExtraC() *Project {
	return executeTemplate(deepCopy(&basicProjectWithExtra), "c")
}

func BasicProjectWithExtraCC() *Project {
	return executeTemplate(deepCopy(&basicProjectWithExtra), "cc")
}

func ComplexProjectC() *Project {
	return executeTemplate(deepCopy(&complexProject), "c")
}

func ComplexProjectCC() *Project {
	return executeTemplate(deepCopy(&complexProject), "cc")
}

func deepCopy(project *Project) *Project {
	bytes, err := yaml.Marshal(project)
	if err != nil {
		panic(err)
	}

	newProject := new(Project)
	if err := yaml.Unmarshal(bytes, newProject); err != nil {
		panic(err)
	}

	return newProject
}

func executeTemplate(project *Project, extension string) *Project {
	project.Name = project.Name + strings.ToUpper(extension)
	for _, node := range project.Nodes {
		node.Name = strings.Replace(
			node.Name,
			".FILE_EXTENSION",
			"."+extension,
			1,
		)

		if strings.HasSuffix(node.Name, ".h") && extension == "cc" {
			class := hex.EncodeToString([]byte(node.Name))
			node.ExtraContent = fmt.Sprintf(`
#ifndef CLASS_%s_H_
#define CLASS_%s_H_

class Class%s {
public:
  Class%s() { foo_ = 1; }

private:
  int foo_;
};

#endif // CLASS_%s_H_
`, class, class, class, class, class)
		}
	}
	return project
}

type ProjectNode struct {
	Name         string   `yaml:"name"`
	Includes     []string `yaml:"includes"`
	Dependencies []string `yaml:"dependencies"`
	ExtraContent string   `yaml:"extra_content"`
}

type Project struct {
	Name  string         `yaml:"name"`
	Root  string         `yaml:"root"`
	Nodes []*ProjectNode `yaml:"nodes"`
}

func (p *Project) String() string {
	bytes, err := yaml.Marshal(p)
	if err != nil {
		return fmt.Sprintf("error marshalling project: %s", err.Error())
	} else {
		return string(bytes)
	}
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
