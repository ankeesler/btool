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
	)
	c.run(
		"./main",
	)
}
