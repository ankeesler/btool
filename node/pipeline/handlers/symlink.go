package handlers

import (
	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/pipeline"
	"github.com/pkg/errors"
)

type symlink struct {
	rf       ResolverFactory
	to, from string
}

// NewSymlink provides a pipeline.Handler that symlinks a file to another file.
func NewSymlink(rf ResolverFactory, to, from string) pipeline.Handler {
	return &symlink{
		rf:   rf,
		to:   to,
		from: from,
	}
}

func (s *symlink) Handle(ctx pipeline.Ctx) error {
	log.Debugf("symlink from %s to %s", s.from, s.to)

	fromN := ctx.Find(s.from)
	if fromN == nil {
		return errors.New("unknown symlink source: " + s.from)
	}

	n := node.New(s.to).Dependency(fromN)
	n.Resolver = s.rf.NewSymlink()
	ctx.Add(n)

	return nil
}

func (s *symlink) String() string { return "symlink" }
