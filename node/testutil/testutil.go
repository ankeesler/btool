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

func DeepCopy(nodes []*node.Node) []*node.Node {
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
