package builder

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ankeesler/btool/builder/compiler"
	"github.com/ankeesler/btool/builder/graph"
	"github.com/ankeesler/btool/builder/linker"
	"github.com/ankeesler/btool/builder/scanner"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Builder struct {
	store string
}

func New(store string) *Builder {
	return &Builder{store: store}
}

func (b *Builder) Build(target string) error {
	if !fileExists(target) {
		return fmt.Errorf("cannot stat target file %s", target)
	}

	g := graph.New()
	s := scanner.New()
	if err := addToGraph(target, g, s); err != nil {
		return errors.Wrap(err, "add to graph")
	}

	nodes, err := g.Sort()
	if err != nil {
		return errors.Wrap(err, "sort graph")
	}

	objects := make([]string, 0, 2)
	c := compiler.New()
	for _, node := range nodes {
		log.Debugf("looking at sorted node %s", node)
		if strings.HasSuffix(node.Name, ".c") {
			log.Debugf("compiling node %s", node)
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

func addToGraph(file string, g *graph.Graph, s *scanner.Scanner) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.Wrap(err, "read file")
	}

	fileNode := &graph.Node{
		Name: file,
	}
	g.Add(fileNode, nil)
	log.Debugf("added %s to graph", fileNode)

	includes, err := s.Scan(data)
	if err != nil {
		return errors.Wrap(err, "scan")
	}
	for _, include := range includes {
		includeNode := &graph.Node{
			Name: include,
		}
		g.Add(fileNode, includeNode)
		log.Debugf("added %s->%s dependency to graph", fileNode, includeNode)

		cFile := strings.ReplaceAll(include, ".h", ".c")
		if fileExists(cFile) {
			cFileNode := &graph.Node{
				Name: cFile,
			}
			g.Add(fileNode, cFileNode)
			log.Debugf("added %s->%s dependency to graph", fileNode, cFileNode)

			if err := addToGraph(cFile, g, s); err != nil {
				return errors.Wrap(err, fmt.Sprintf("add %s to graph", cFile))
			}
		}
	}

	return nil
}

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	} else {
		return false
	}
}
