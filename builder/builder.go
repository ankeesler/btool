package builder

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ankeesler/btool/builder/compiler"
	"github.com/ankeesler/btool/builder/graph"
	"github.com/ankeesler/btool/builder/linker"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type Builder struct {
	fs    afero.Fs
	store string
}

func New(fs afero.Fs, store string) *Builder {
	return &Builder{
		fs:    fs,
		store: store,
	}
}

func (b *Builder) Build(graph *graph.Graph) error {
	logrus.Info("building graph")

	nodes, err := graph.Sort()
	if err != nil {
		return errors.Wrap(err, "sort graph")
	}

	objects := make([]string, 0, 2)
	c := compiler.New()
	for _, node := range nodes {
		logrus.Debugf("looking at sorted node %s", node)
		if strings.HasSuffix(node.Name, ".c") {
			logrus.Debugf("compiling node %s", node)
			output := filepath.Join(
				"/tmp",
				strings.ReplaceAll(filepath.Base(node.Name), ".c", ".o"),
			)
			if err := c.Compile(node.Name, output); err != nil {
				return errors.Wrap(err, fmt.Sprintf("compile %s", node.Name))
			}
			objects = append(objects, output)
		}
	}

	l := linker.New()
	output := filepath.Join(
		b.store,
		strings.ReplaceAll(
			filepath.Base(target),
			".c",
			"",
		),
	)
	if err := l.Link(output, objects); err != nil {
		return errors.Wrap(err, fmt.Sprintf("link %s", output))
	}

	return nil
}
