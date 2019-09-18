// Package runner provides functionality to run an executable node.Node.
package runner

import (
	"os"
	"os/exec"
	"strings"

	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Callback

// Callback is an interface that clients can use to be notified when a node.Node
// is run.
type Callback interface {
	OnRun(*node.Node)
}

// Runner is a type that can run a node.Node.
type Runner struct {
	c Callback
}

// New creates a new Runner.
func New(c Callback) *Runner {
	return &Runner{
		c: c,
	}
}

func (r *Runner) Run(n *node.Node) error {
	r.c.OnRun(n)

	// If we get a node.Node named something like 'main', we want to actually run
	// './main'.
	name := n.Name
	if strings.Index(n.Name, "/") == -1 {
		name = "./" + name
	}

	cmd := exec.Command(name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "run")
	}

	return nil
}
