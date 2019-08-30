package integration

import (
	"os"
	"path/filepath"

	"github.com/stretchr/testify/require"
)

func googletest(c *config) {
	here, err := os.Getwd()
	require.Nil(c.t, err)
	c.run(
		c.btool,
		"-target",
		filepath.Join(c.cache, "projects", "googletest", "gtest.a"),
		"-root",
		c.root,
		"-cache",
		c.cache,
		"-loglevel",
		"debug",
		"-registry",
		filepath.Join(here, "..", "data"),
		"-output",
		"gtest.a",
	)
	c.run(
		"ls",
		"gtest.a",
	)
}

func googletestTest(c *config) {
	c.run(
		c.btool,
		"-target",
		"dep-1/dep-1-test",
		"-root",
		c.root,
		"-cache",
		c.cache,
		"-loglevel",
		"debug",
	)
	c.run(
		"dep-1/dep-1-test",
	)
}
