package integration_test

import (
	"os/exec"
	"testing"
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
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = c.wd
	if output, err := cmd.CombinedOutput(); err != nil {
		c.t.Fatal(err, ":", string(output))
	}
}

type testCase struct {
	name     string
	testFunc func(c *config)
}
