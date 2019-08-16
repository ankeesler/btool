package integration

import (
	"os"
	"path/filepath"

	"github.com/stretchr/testify/require"
)

func executable(c *config) {
	c.run(
		c.btool,
		"-target",
		"main",
		"-root",
		c.root,
		"-cache",
		c.cache,
		"-output",
		"main",
		"-loglevel",
		"debug",
	)
	c.run(
		"./main",
	)
}

func executableLocalRegistry(c *config) {
	dir, err := os.Getwd()
	require.Nil(c.t, err)
	registryData := filepath.Join(dir, "..", "data")

	c.run(
		c.btool,
		"-target",
		"main",
		"-root",
		c.root,
		"-cache",
		c.cache,
		"-registries",
		registryData,
		"-output",
		"main",
		"-loglevel",
		"debug",
	)
	c.run(
		"./main",
	)
}

func executableRunTwice(c *config) {
	for i := 0; i < 2; i++ {
		c.run(
			c.btool,
			"-target",
			"main",
			"-root",
			c.root,
			"-cache",
			c.cache,
			"-output",
			"main",
			"-loglevel",
			"debug",
		)
		c.run(
			"./main",
		)
	}
}
