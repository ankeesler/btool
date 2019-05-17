package testutil

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/ankeesler/btool/scanner/graph"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

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
