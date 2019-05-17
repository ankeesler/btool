// Package scanner provides the ability to build a dependency graph for a C/C++
// project.
package scanner

import (
	"github.com/ankeesler/btool/builder/graph"
	"github.com/spf13/afero"
)

type Scanner struct {
	fs afero.Fs
}

func New(fs afero.Fs) *Scanner {
	return &Scanner{
		fs: fs,
	}
}

func (s *Scanner) Scan(root string) (*graph.Graph, error) {
	return nil, nil
}
