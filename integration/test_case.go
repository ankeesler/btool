package integration

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

type config struct {
	btool   string
	root    string
	cache   string
	wd      string
	fixture string

	t *testing.T
}

func (c *config) run(args ...string) {
	c.t.Log("run", args)

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = c.wd
	output, err := cmd.CombinedOutput()
	c.t.Log("output", string(output))
	require.Nil(c.t, err)
}

type testCase struct {
	name     string
	testFunc func(c *config)
}
