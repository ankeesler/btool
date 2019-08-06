package integration

import "path/filepath"

func executable(c *config) {
	c.run(
		c.btool,
		"-target",
		"main",
		"-root",
		c.root,
		"-cache",
		c.cache,
	)
	c.run(
		filepath.Join(c.cache, filepath.Base(c.root), "main"),
	)
}
