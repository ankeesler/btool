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
	)
}
