package testutil

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/ankeesler/btool/node"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type Nodes []*node.Node

func (nodes Nodes) WithoutDependencies() Nodes {
	for _, n := range nodes {
		n.Dependencies = nil
	}
	return nodes
}

func (nodes Nodes) WithObjects() Nodes {
	for _, n := range nodes {
		n.Objects = make([]string, 0)
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

func (nodes Nodes) PopulateFS(fs afero.Fs) {
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
