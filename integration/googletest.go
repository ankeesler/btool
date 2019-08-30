package integration

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
		"dep-1/dep-1-test",
	)
}
