package integration

import "path/filepath"

func googletest(c *config) {
	gtesta := filepath.Join(
		c.cache,
		"googletest",
		"library",
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
		"http://127.0.0.1:8080",
	)
	c.run(
		"ls",
		gtesta,
	)
}
