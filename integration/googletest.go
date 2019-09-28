package integration

func googletest(c *config) {
	if c.example != "BasicCC" {
		return
	}

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
