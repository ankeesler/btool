// Package compiler provides a node.Handler that compiles each node into an object
// archive.
package compiler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . C

// C is an actual compiler client.
type C interface {
	CompileC(output, input string, includeDirs []string) error
	CompileCC(output, input string, includeDirs []string) error
}

type Compiler struct {
	c     C
	fs    afero.Fs
	root  string
	cache string
}

func New(c C, fs afero.Fs, root, cache string) *Compiler {
	return &Compiler{
		c:     c,
		fs:    fs,
		root:  root,
		cache: cache,
	}
}

func (c *Compiler) Handle(nodes []*node.Node) ([]*node.Node, error) {
	for _, node := range nodes {
		if err := c.handleNode(node); err != nil {
			return nil, errors.Wrap(err, "handle node "+node.Name)
		}
	}
	return nodes, nil
}

func (c *Compiler) handleNode(n *node.Node) error {
	if !c.needsCompile(n) {
		return nil
	}

	n.Objects = make([]string, 0)
	for _, source := range n.Sources {
		object, err := c.handleSource(n, source)
		if err != nil {
			return errors.Wrap(err, "handle source "+source)
		}

		n.Objects = append(n.Objects, object)
	}
	return nil
}

func (c *Compiler) needsCompile(n *node.Node) bool {
	return true
}

func (c *Compiler) handleSource(n *node.Node, source string) (string, error) {
	var compileFunc func(string, string, []string) error
	if strings.HasSuffix(source, ".c") {
		compileFunc = c.c.CompileC
	} else if strings.HasSuffix(n.Name, ".cc") {
		compileFunc = c.c.CompileCC
	} else {
		return "", fmt.Errorf("file is not compilable: %s", source)
	}

	inputFile := filepath.Join(
		c.root,
		source,
	)

	outputFile := filepath.Join(
		c.objectsDir(),
		filepath.Dir(source),
		filepath.Base(source)+".o",
	)

	if older, err := c.isFileOlder(
		outputFile,
		append(n.Dependencies, n),
	); err != nil {
		return "", errors.Wrap(err, "is file older")
	} else if older {
		logrus.Infof("compile: %s (up to date)", n.Name)
		return outputFile, nil
	} else {
		logrus.Infof("compile: %s", n.Name)
	}

	outputDir := filepath.Dir(outputFile)
	if err := c.fs.MkdirAll(outputDir, 0700); err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("mkdir (%s)", outputDir))
	}

	if err := compileFunc(
		outputFile,
		inputFile,
		append(n.IncludePaths, c.root),
	); err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("compile %s", outputFile))
	}

	return outputFile, nil
}

func (c *Compiler) isFileOlder(from string, tos []*node.Node) (bool, error) {
	logrus.Debugf("is %s older than %s", from, tos)

	fromStat, err := c.fs.Stat(from)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.Debugf("%s does not exist", from)
			return false, nil
		} else {
			return false, errors.Wrap(err, "stat from")
		}
	}

	for _, to := range tos {
		toStat, err := c.fs.Stat(filepath.Join(c.root, to.Name))
		if err != nil {
			if os.IsNotExist(err) {
				logrus.Debugf("%s does not exist", to)
				return false, nil
			} else {
				return false, errors.Wrap(err, "stat to")
			}
		}

		fromModTime := fromStat.ModTime()
		toModTime := toStat.ModTime()
		if fromModTime.Before(toModTime) {
			logrus.Debugf(
				"%s (%s) is not older than %s (%s)",
				from,
				fromModTime,
				to,
				toModTime,
			)
			return false, nil
		}
	}

	return true, nil
}

func (c *Compiler) objectsDir() string {
	return filepath.Join(c.cache, "objects")
}
