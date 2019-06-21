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
	yaml "gopkg.in/yaml.v2"
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
				"<stdio.h>",
				"\"master.h\"",
				"\"dep-0/dep-0a.h\"",
			},
			Dependencies: []string{
				"master.h",
				"dep-0/dep-0a.h",
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
			Name:         "dep-1/dep-1a.h",
			Includes:     []string{},
			Dependencies: []string{},
		},
		//&ProjectNode{
		//	Name: "dep-1/dep-1a-test.FILE_EXTENSION",
		//	Includes: []string{
		//		"\"dep-1a.h\"",
		//		"\"gtest/gtest.h\"",
		//	},
		//	Dependencies: []string{
		//		"dep-1/dep-1a.h",
		//	},
		//},
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

func BigProjectC() *Project {
	return generateBigProject("c")
}

func BigProjectCC() *Project {
	return generateBigProject("cc")
}

func generateBigProject(extension string) *Project {
	p := &Project{
		Name:  "BigProject" + strings.ToUpper(extension),
		Root:  "/tuna/root",
		Nodes: make([]*ProjectNode, 0),
	}

	// Generate a "common" directory with 100 .h files
	// No dependencies in each file.
	for i := 0; i < 100; i++ {
		prefix := fmt.Sprintf("common/common-%02d.", i)
		includeNode := &ProjectNode{
			Name:         prefix + "h",
			Includes:     []string{},
			Dependencies: []string{},
		}
		if extension == "cc" {
			includeNode.ExtraContent = generateClass(prefix + "h")
		}

		p.Nodes = append(p.Nodes, includeNode)
		p.Nodes = append(p.Nodes, &ProjectNode{
			Name: prefix + extension,
			Includes: []string{
				"\"" + prefix + "h" + "\"",
			},
			Dependencies: []string{
				prefix + "h",
			},
		})
	}

	// Generate 25 .h files for 4 levels.
	// Give each .h file 4 "common" dependencies.
	for i := 0; i < 100; i++ {
		var prefix string
		for j := 0; j < (i/25)+1; j++ {
			prefix += fmt.Sprintf("level-%02d/", j)
		}

		prefix += fmt.Sprintf("file-%02d.", i)
		includeNode := &ProjectNode{
			Name:         prefix + "h",
			Includes:     []string{},
			Dependencies: []string{},
		}
		if extension == "cc" {
			includeNode.ExtraContent = generateClass(prefix + "h")
		}

		sourceNode := &ProjectNode{
			Name: prefix + extension,
			Includes: []string{
				"\"" + prefix + "h" + "\"",
			},
			Dependencies: []string{
				prefix + "h",
			},
		}

		for j := 0; j < 4; j++ {
			include := fmt.Sprintf("common/common-%02d.h", (i*j)%100)
			includeQuoted := fmt.Sprintf("\"%s\"", include)
			includeNode.Includes = append(includeNode.Includes, includeQuoted)
			includeNode.Dependencies = append(includeNode.Dependencies, include)
			sourceNode.Includes = append(sourceNode.Includes, includeQuoted)
			sourceNode.Dependencies = append(sourceNode.Dependencies, include)
		}

		p.Nodes = append(p.Nodes, includeNode)
		p.Nodes = append(p.Nodes, sourceNode)
	}

	// Generate main with 10 includes from the common directory
	// and 10 includes from the level-00 directory.
	mainNode := &ProjectNode{
		Name: "main." + extension,
		Includes: []string{
			"<stdio.h>",
		},
		Dependencies: []string{},
		ExtraContent: `
int main(int argc, char *argv[]) {
	printf("hey! i am a large project!\n");
	return 0;
}
`,
	}
	for i := 0; i < 10; i++ {
		include := fmt.Sprintf("common/common-%02d.h", i)
		includeQuoted := fmt.Sprintf("\"%s\"", include)
		mainNode.Includes = append(mainNode.Includes, includeQuoted)
		mainNode.Dependencies = append(mainNode.Dependencies, include)

		include = fmt.Sprintf("level-00/file-%02d.h", i)
		includeQuoted = fmt.Sprintf("\"%s\"", include)
		mainNode.Includes = append(mainNode.Includes, includeQuoted)
		mainNode.Dependencies = append(mainNode.Dependencies, include)
	}

	p.Nodes = append(p.Nodes, mainNode)

	return p
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
			node.ExtraContent = generateClass(node.Name)
		}
	}
	return project
}

func generateClass(name string) string {
	class := hex.EncodeToString([]byte(name))
	return fmt.Sprintf(`
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
			dependencyPath := dependency
			if !filepath.IsAbs(dependencyPath) {
				dependencyPath = filepath.Join(p.Root, dependency)
			}
			dependencyNode := &graph.Node{
				Name: dependencyPath,
			}
			g.Add(pathNode, dependencyNode)
		}
	}

	return g
}
