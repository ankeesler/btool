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
		"-debug",
	)
	c.run(
		filepath.Join(c.root, "main"),
	)

	c.run(
		c.btool,
		"-target",
		"main",
		"-root",
		c.root,
		"-cache",
		c.cache,
		"-debug",
		"-clean",
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
		"-registry",
		registryData,
		"-debug",
	)
	c.run(
		filepath.Join(c.root, "main"),
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
		"-debug",
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
			"-debug",
		)
		c.run(
			filepath.Join(c.root, "main"),
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
		"-debug",
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
		"-debug",
		"-cache",
		cache,
	)
	c.run(
		filepath.Join(c.root, "main"),
	)

	c.run(
		c.btool,
		"-target",
		"main",
		"-debug",
		"-cache",
		cache,
		"-clean",
	)
}
