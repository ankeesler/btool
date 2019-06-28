// Package testutil provides common node test utilities.
package testutil

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/ankeesler/btool/node"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

var (
	Dep0h = node.Node{
		Name:         "dep-0/dep-0.h",
		Sources:      []string{},
		Headers:      []string{"dep-0/dep-0.h"},
		Dependencies: []*node.Node{},
	}
	Dep0c = node.Node{
		Name:         "dep-0/dep-0.c",
		Sources:      []string{"dep-0/dep-0.c"},
		Headers:      []string{},
		Dependencies: []*node.Node{&Dep0h},
	}

	Dep1h = node.Node{
		Name:         "dep-1/dep-1.h",
		Sources:      []string{},
		Headers:      []string{"dep-1/dep-1.h"},
		Dependencies: []*node.Node{&Dep0h},
	}
	Dep1c = node.Node{
		Name:         "dep-1/dep-1.c",
		Sources:      []string{"dep-1/dep-1.c"},
		Headers:      []string{},
		Dependencies: []*node.Node{&Dep1h, &Dep0h},
	}

	Mainc = node.Node{
		Name:         "main.c",
		Sources:      []string{"main.c"},
		Headers:      []string{},
		Dependencies: []*node.Node{&Dep1h, &Dep0h},
	}
)

var (
	BasicNodes = []*node.Node{
		&Dep0c,
		&Dep0h,
		&Dep1c,
		&Dep1h,
		&Mainc,
	}
)

func RemoveDependencies(nodes []*node.Node) []*node.Node {
	newNodes := make([]*node.Node, len(nodes))
	copy(newNodes, nodes)
	for _, n := range newNodes {
		n.Dependencies = nil
	}
	return newNodes
}

func PopulateFS(nodes []*node.Node, fs afero.Fs) {
	for _, node := range nodes {
		content := bytes.NewBuffer([]byte{})
		content.WriteString(fmt.Sprintf("// %s\n", node.Name))
		for _, dependency := range node.Dependencies {
			for _, header := range dependency.Headers {
				content.WriteString(fmt.Sprintf("\n#include \"%s\"", header))
			}
		}
		//if node.ExtraContent != "" {
		//	content.WriteString(node.ExtraContent)
		//}

		file := filepath.Join("/", node.Name)
		dir := filepath.Dir(file)
		if err := fs.MkdirAll(dir, 0700); err != nil {
			panic(err)
		}
		if err := afero.WriteFile(fs, file, content.Bytes(), 0700); err != nil {
			panic(err)
		}

		logrus.Debugf("created file " + file)
	}
}
