package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/pipeline"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type resolve struct {
	fs afero.Fs
}

// NewResolve returns a pipeline.Handler that runs all of the node.Resolver graph
// for a particular target.
func NewResolve(fs afero.Fs) pipeline.Handler {
	return &resolve{
		fs: fs,
	}
}

func (r *resolve) Handle(ctx *pipeline.Ctx) {
	target := ctx.KV[pipeline.CtxTarget]
	nodes := ctx.Nodes

	n := node.Find(target, nodes)
	if n == nil {
		ctx.Err = fmt.Errorf("unknown target %s", target)
		return
	}

	built := make(map[*node.Node]bool)
	if _, err := r.resolve(ctx, n, built); err != nil {
		ctx.Err = errors.Wrap(err, "build")
		return
	}
}

func (r *resolve) Name() string { return "resolve" }

func (r *resolve) resolve(
	ctx *pipeline.Ctx,
	n *node.Node,
	built map[*node.Node]bool,
) (time.Time, error) {
	if built[n] {
		return time.Time{}, nil
	}

	logrus.Debugf("building %s", n.Name)

	latestT := time.Unix(0, 0)
	for _, d := range n.Dependencies {
		logrus.Debugf("building dependency %s", d.Name)

		t, err := r.resolve(ctx, d, built)
		if err != nil {
			return time.Time{}, errors.Wrap(err, "build "+d.Name)
		}

		if t.After(latestT) {
			latestT = t
		}
	}

	file := filepath.Join(ctx.KV[pipeline.CtxRoot], n.Name)
	logrus.Debugf("resolving %s", file)

	var t time.Time
	if n.Resolver != nil {
		exists := true
		stat, err := r.fs.Stat(file)
		if err != nil {
			if !os.IsNotExist(err) {
				return time.Time{}, errors.Wrap(err, "stat "+file)
			} else {
				exists = false
			}
		}

		if !exists || latestT.After(stat.ModTime()) {
			if err := n.Resolver.Resolve(n); err != nil {
				return time.Time{}, errors.Wrap(err, "resolve "+n.Name)
			}
		}

		stat, err = r.fs.Stat(file)
		if err != nil {
			return time.Time{}, errors.Wrap(err, "stat "+file)
		}
		t = stat.ModTime()
	}

	built[n] = true

	return t, nil
}
