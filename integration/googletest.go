package integration

import "path/filepath"

func googletest(c *config) {
	gtesta := filepath.Join(
		c.cache,
		"projects",
		"googletest",
		"gtest.a",
	)
	c.run(
		c.btool,
		"-target",
		gtesta,
		"-root",
		c.root,
		"-cache",
		c.cache,
		"-loglevel",
		"debug",
		"-registries",
		"/Users/ankeesler/workspace/btool/data",
	)
	c.run(
		"ls",
		gtesta,
	)
}
