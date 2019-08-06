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
		ctx.Err = errors.Wrap(err, "resolve")
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

	logrus.Debugf("resolving %s", n.Name)

	latestT := time.Unix(0, 0)
	for _, d := range n.Dependencies {
		logrus.Debugf("resolving dependency %s", d.Name)

		t, err := r.resolve(ctx, d, built)
		if err != nil {
			return time.Time{}, errors.Wrap(err, "resolve "+d.Name)
		}

		if t.After(latestT) {
			latestT = t
		}
	}

	var t time.Time
	exists := true
	stat, err := r.fs.Stat(n.Name)
	if err != nil {
		if !os.IsNotExist(err) {
			return time.Time{}, errors.Wrap(err, "stat "+n.Name)
		} else {
			exists = false
		}
	} else {
		t = stat.ModTime()
	}

	if n.Resolver != nil && (!exists || latestT.After(stat.ModTime())) {
		logrus.Debugf("really resolving %s", n.Name)

		dir := filepath.Dir(n.Name)
		if err := r.fs.MkdirAll(dir, 0755); err != nil {
			return time.Time{}, errors.Wrap(err, "mkdir "+dir)
		}

		if err := n.Resolver.Resolve(n); err != nil {
			return time.Time{}, errors.Wrap(err, "really resolve "+n.Name)
		}

		stat, err = r.fs.Stat(n.Name)
		if err != nil {
			return time.Time{}, errors.Wrap(err, "stat "+n.Name)
		}
		t = stat.ModTime()
	}

	built[n] = true

	return t, nil
}
