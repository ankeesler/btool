// Package scanner provides the ability to build a dependency graph for a C/C++
// project.
package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ankeesler/btool/config"
	"github.com/ankeesler/btool/scanner/graph"
	"github.com/ankeesler/btool/scanner/includes"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type Scanner struct {
	fs     afero.Fs
	config *config.Config

	graph *graph.Graph
}

func New(fs afero.Fs, config *config.Config) *Scanner {
	return &Scanner{
		fs:     fs,
		config: config,
	}
}

func (s *Scanner) ScanFile(file string) (*graph.Graph, error) {
	logrus.Info("scanning from file " + file)

	s.graph = graph.New()

	logrus.Debugf("walking dependencies from file %s", file)
	if err := s.walk(file, make(map[string]bool)); err != nil {
		return nil, errors.Wrap(err, "walk")
	}

	return s.graph, nil
}

func (s *Scanner) ScanRoot() (*graph.Graph, error) {
	logrus.Info("scanning from root " + s.config.Root)

	s.graph = graph.New()

	logrus.Debugf("walking fs from root %s", s.config.Root)
	if err := afero.Walk(
		s.fs,
		s.config.Root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return errors.Wrap(err, "walk")
			}

			if info.IsDir() {
				logrus.Debugf("skipping directory %s", path)
				return nil
			}

			if err := s.addToGraph(path); err != nil {
				return errors.Wrap(err, fmt.Sprintf("add %s to graph", path))
			}

			return nil
		},
	); err != nil {
		return nil, errors.Wrap(err, "walk")
	}

	return s.graph, nil
}

// Everything in the graph is absolute! It is easier that way since the path
// we are passed in this function is absolute.
func (s *Scanner) addToGraph(path string) error {
	data, err := afero.ReadFile(s.fs, path)
	if err != nil {
		return errors.Wrap(err, "read file "+path)
	}
	logrus.Debugf("read file %s", path)

	pathNode := &graph.Node{
		Name: path,
	}
	s.graph.Add(pathNode, nil)

	includes, err := includes.Parse(data)
	if err != nil {
		return errors.Wrap(err, "parse includes")
	}
	logrus.Debugf("parsed includes %s", includes)

	for _, include := range includes {
		includePath, err := s.resolveIncludePath(include, filepath.Dir(path))
		if err != nil {
			return errors.Wrap(err, "resolve include path "+include)
		}

		includeNode := &graph.Node{
			Name: includePath,
		}
		s.graph.Add(pathNode, includeNode)
	}

	return nil
}

func (s *Scanner) resolveIncludePath(include, dir string) (string, error) {
	rootRelJoin := filepath.Join(s.config.Root, include)
	if s.exists(rootRelJoin) {
		return rootRelJoin, nil
	}

	dirRelJoin := filepath.Join(dir, include)
	if s.exists(dirRelJoin) {
		return filepath.Join(dir, include), nil
	}

	return "", errors.New("cannot resolve include: " + include)
}

func (s *Scanner) walk(file string, visited map[string]bool) error {
	if visited[file] {
		return nil
	}

	visited[file] = true

	data, err := afero.ReadFile(s.fs, file)
	if err != nil {
		return errors.Wrap(err, "read file "+file)
	}
	logrus.Debugf("read file %s", file)

	fileNode := &graph.Node{
		Name: file,
	}
	s.graph.Add(fileNode, nil)

	includes, err := includes.Parse(data)
	if err != nil {
		return errors.Wrap(err, "parse includes")
	}
	logrus.Debugf("parsed includes %s", includes)

	for _, include := range includes {
		includePath, err := s.resolveIncludePath(include, filepath.Dir(file))
		if err != nil {
			return errors.Wrap(err, "resolve include path "+include)
		}

		includeNode := &graph.Node{
			Name: includePath,
		}
		s.graph.Add(fileNode, includeNode)

		sources := s.sourcesForInclude(includePath)
		for _, source := range sources {
			if err := s.walk(source, visited); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Scanner) sourcesForInclude(includePath string) []string {
	sources := make([]string, 0)

	sourcePathC := strings.Replace(includePath, ".h", ".c", 1)
	if s.exists(sourcePathC) {
		sources = append(sources, sourcePathC)
	}

	sourcePathCC := strings.Replace(includePath, ".h", ".cc", 1)
	if s.exists(sourcePathCC) {
		sources = append(sources, sourcePathCC)
	}

	return sources
}

func (s *Scanner) exists(path string) bool {
	exists, _ := afero.Exists(s.fs, path)
	return exists
}
