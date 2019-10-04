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
		"-loglevel",
		"debug",
		"-run",
	)

	c.run(
		c.btool,
		"-target",
		"main",
		"-root",
		c.root,
		"-cache",
		c.cache,
		"-loglevel",
		"debug",
		"-clean",
	)
}

func executableLocalRegistry(c *config) {
	dir, err := os.Getwd()
	require.Nil(c.t, err)
	registryData := filepath.Join(dir, "..", "registry", "data")

	c.run(
		c.btool,
		"-target",
		"main",
		"-root",
		c.root,
		"-cache",
		c.cache,
		"-registry",
		registryData,
		"-loglevel",
		"debug",
		"-run",
	)

	c.run(
		c.btool,
		"-target",
		"main",
		"-root",
		c.root,
		"-cache",
		c.cache,
		"-registry",
		registryData,
		"-loglevel",
		"debug",
		"-clean",
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
			"-loglevel",
			"debug",
			"-run",
		)
	}

	c.run(
		c.btool,
		"-target",
		"main",
		"-root",
		c.root,
		"-cache",
		c.cache,
		"-loglevel",
		"debug",
		"-clean",
	)
}

func executableSubdirectoryCache(c *config) {
	wd, err := os.Getwd()
	require.Nil(c.t, err)

	c.wd = filepath.Join(wd, "..", "example", c.example)
	cache := filepath.Join(c.wd, "cache")
	defer os.RemoveAll(cache)

	c.run(
		c.btool,
		"-target",
		"main",
		"-loglevel",
		"debug",
		"-cache",
		cache,
		"-run",
	)

	c.run(
		c.btool,
		"-target",
		"main",
		"-loglevel",
		"debug",
		"-cache",
		cache,
		"-clean",
	)
}
