package resolvers

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
)

type script struct {
	s string
}

// NewScript returns a node.Resolver that runs a script.
func NewScript(s string) node.Resolver {
	return &script{
		s: s,
	}
}

func (s *script) Resolve(n *node.Node) error {
	cmd := exec.Command("bash")
	cmd.Args = append(cmd.Args, "-c")
	cmd.Args = append(cmd.Args, s.s)

	stdout := bytes.NewBuffer([]byte{})
	stderr := bytes.NewBuffer([]byte{})
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	log.Debugf("script: running %s", cmd.Args)

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, stderr.String())
	}

	if err := os.MkdirAll(filepath.Dir(n.Name), 0755); err != nil {
		return errors.Wrap(err, "mkdir all")
	}

	if err := ioutil.WriteFile(n.Name, stdout.Bytes(), 0644); err != nil {
		return errors.Wrap(err, "write file")
	}

	return nil
}
