package testutil

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ankeesler/btool/node"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type Nodes []*node.Node

func (nodes Nodes) WithoutDependencies() Nodes {
	for _, n := range nodes {
		n.Dependencies = make([]*node.Node, 0)
	}
	return nodes
}

// Copy performs a deep copy.
func (nodes Nodes) Copy() Nodes {
	oldNew := make(map[*node.Node]*node.Node)
	newNodes := make([]*node.Node, 0, len(nodes))

	for _, n := range nodes {
		newNode := new(node.Node)
		*newNode = *n
		newNode.Dependencies = []*node.Node{}
		oldNew[n] = newNode

		newNodes = append(newNodes, newNode)
	}

	for _, n := range nodes {
		newNode := oldNew[n]
		for _, d := range n.Dependencies {
			newNode.Dependencies = append(newNode.Dependencies, oldNew[d])
		}
	}

	return newNodes
}

// Cast is a utility function to cast a Nodes to a []*node.Node.
func (nodes Nodes) Cast() []*node.Node {
	return []*node.Node(nodes)
}

func (nodes Nodes) PopulateFS(root string, fs afero.Fs) {
	for _, node := range nodes {
		content := bytes.NewBuffer([]byte{})
		content.WriteString(fmt.Sprintf("// %s\n", node.Name))
		for _, dependency := range node.Dependencies {
			if strings.HasSuffix(dependency.Name, ".h") {
				content.WriteString(fmt.Sprintf("\n#include \"%s\"", dependency.Name))
			}
		}

		if node.Name == "main.c" {
			content.WriteString("\nint main(int argc, char *argv[]) { return 0; }")
		}

		file := filepath.Join(root, node.Name)
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
