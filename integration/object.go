package integration

func object(c *config) {
	c.run(
		c.btool,
		"-target",
		"dep-0/dep-0.o",
		"-root",
		c.root,
		"-cache",
		c.cache,
		"-loglevel",
		"debug",
		"-output",
		"dep-0/dep-0.o",
	)
	c.run(
		"ls",
		"dep-0/dep-0.o",
	)
}
