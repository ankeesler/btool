package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ankeesler/btool/scanner/graph"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func (b *Builder) Build(graph *graph.Graph) error {
	logrus.Info("building graph")

	nodes, err := graph.Sort()
	if err != nil {
		return errors.Wrap(err, "sort graph")
	}

	objects := make([]string, 0)

	for _, node := range nodes {
		logrus.Debugf("looking at sorted node %s", node)

		var object string
		var err error
		if strings.HasSuffix(node.Name, ".c") {
			object, err = b.compile(node, ".c", b.c.CompileC)
		} else if strings.HasSuffix(node.Name, ".cc") {
			object, err = b.compile(node, ".cc", b.c.CompileCC)
		}

		if err != nil {
			return errors.Wrap(err, "compile")
		}

		if object != "" {
			objects = append(objects, object)
		}
	}

	if err := b.link(objects); err != nil {
		return errors.Wrap(err, "link")
	}

	return nil
}

func (b *Builder) compile(
	node *graph.Node,
	extension string,
	compileFunc func(output, input, include string) error,
) (string, error) {
	rootRelNodeName, err := filepath.Rel(b.root, node.Name)
	if err != nil {
		return "", errors.Wrap(err, "rel")
	}

	outputFile := filepath.Join(
		b.store,
		"objects",
		strings.Replace(rootRelNodeName, extension, ".o", 1),
	)

	if older, err := b.isFileOlder(outputFile, node.Name); err != nil {
		return "", errors.Wrap(err, "is file older")
	} else if older {
		logrus.Infof("compile: %s (up to date)", node.Name)
		return outputFile, nil
	} else {
		logrus.Infof("compile: %s", node.Name)
	}

	outputDir := filepath.Dir(outputFile)
	if err := b.fs.MkdirAll(outputDir, 0700); err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("mkdir (%s)", outputDir))
	}

	if err := compileFunc(outputFile, node.Name, b.root); err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("compile %s", outputFile))
	}

	return outputFile, nil
}

func (b *Builder) link(objects []string) error {
	outputFile := filepath.Join(
		b.store,
		"binaries",
		"out",
	)

	if older, err := b.isFileOlder(outputFile, objects...); err != nil {
		return errors.Wrap(err, "is file older")
	} else if older {
		logrus.Infof("link: %s (up to date)", outputFile)
		return nil
	} else {
		logrus.Infof("link: %s", outputFile)
	}

	outputDir := filepath.Dir(outputFile)
	if err := b.fs.MkdirAll(outputDir, 0700); err != nil {
		return errors.Wrap(err, fmt.Sprintf("mkdir (%s)", outputDir))
	}

	if err := b.l.Link(outputFile, objects); err != nil {
		return errors.Wrap(err, fmt.Sprintf("link %s", outputFile))
	}

	return nil
}

func (b *Builder) isFileOlder(from string, tos ...string) (bool, error) {
	logrus.Debugf("is %s older than %s", from, tos)

	fromStat, err := b.fs.Stat(from)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, errors.Wrap(err, "stat from")
		}
	}

	for _, to := range tos {
		toStat, err := b.fs.Stat(to)
		if err != nil {
			if os.IsNotExist(err) {
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

func convertFilesToNames(files []afero.File) []string {
	names := make([]string, 0, len(files))
	for _, file := range files {
		names = append(names, file.Name())
	}
	return names
}
