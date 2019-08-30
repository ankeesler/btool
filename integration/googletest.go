package integration

import "path/filepath"

func googletest(c *config) {
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
		filepath.Join(c.root, "dep-1/dep-1-test"),
	)
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
		"-clean",
	)
}
