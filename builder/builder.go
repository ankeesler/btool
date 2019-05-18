// Package builder is in charge of walking the project graph to compile and link
// stuff.
package builder

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ankeesler/btool/builder/compiler"
	"github.com/ankeesler/btool/builder/linker"
	"github.com/ankeesler/btool/scanner/graph"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type Builder struct {
	fs    afero.Fs
	root  string
	store string
	c     *compiler.Compiler
	l     *linker.Linker
}

func New(fs afero.Fs, root, store string, c *compiler.Compiler, l *linker.Linker) *Builder {
	return &Builder{
		fs:    fs,
		root:  root,
		store: store,
		c:     c,
		l:     l,
	}
}

func (b *Builder) Build(graph *graph.Graph) error {
	logrus.Info("building graph")

	nodes, err := graph.Sort()
	if err != nil {
		return errors.Wrap(err, "sort graph")
	}

	objects := make([]afero.File, 0, 2)
	defer func() {
		for _, object := range objects {
			object.Close()
		}
	}()

	for _, node := range nodes {
		logrus.Debugf("looking at sorted node %s", node)
		if strings.HasSuffix(node.Name, ".c") {
			object, err := b.compile(node)
			if err != nil {
				return errors.Wrap(err, "compile")
			}
			objects = append(objects, object)
		}
	}

	if err := b.link(objects); err != nil {
		return errors.Wrap(err, "link")
	}

	return nil
}

func (b *Builder) compile(node *graph.Node) (afero.File, error) {
	logrus.Infof("compiling node %s", node)

	rootRelNodeName, err := filepath.Rel(b.root, node.Name)
	if err != nil {
		return nil, errors.Wrap(err, "rel")
	}

	outputFile := filepath.Join(
		b.store,
		"objects",
		strings.Replace(rootRelNodeName, ".c", ".o", 1),
	)
	outputDir := filepath.Dir(outputFile)
	if err := b.fs.MkdirAll(outputDir, 0700); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("mkdir (%s)", outputDir))
	}

	output, err := b.fs.Create(outputFile)
	if err != nil {
		return nil, errors.Wrap(err, "create output")
	}
	// should be closed by the caller!

	input, err := b.fs.Open(node.Name)
	if err != nil {
		return nil, errors.Wrap(err, "open input")
	}
	defer input.Close()

	if err := b.c.Compile(output.Name(), input.Name(), b.root); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("compile %s", output.Name()))
	}

	return output, nil
}

func (b *Builder) link(objects []afero.File) error {
	logrus.Infof("linking")

	outputFile := filepath.Join(
		b.store,
		"binaries",
		"out",
	)
	outputDir := filepath.Dir(outputFile)
	if err := b.fs.MkdirAll(outputDir, 0700); err != nil {
		return errors.Wrap(err, fmt.Sprintf("mkdir (%s)", outputDir))
	}

	output, err := b.fs.Create(outputFile)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("create output (%s)", outputFile))
	}
	defer output.Close()

	if err := b.l.Link(output.Name(), convertFilesToNames(objects)); err != nil {
		return errors.Wrap(err, fmt.Sprintf("link %s", output.Name()))
	}

	return nil
}

func convertFilesToNames(files []afero.File) []string {
	names := make([]string, 0, len(files))
	for _, file := range files {
		names = append(names, file.Name())
	}
	return names
}
