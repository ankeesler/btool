package integration

func object(c *config) {
	for i := 0; i < 2; i++ {
		c.run(
			c.btool,
			"-target",
			"dep-0/dep-0.o",
			"-root",
			c.root,
			"-cache",
			c.cache,
			"-debug",
		)
		c.run(
			"ls",
			"dep-0/dep-0.o",
		)
	}

	c.run(
		c.btool,
		"-target",
		"dep-0/dep-0.o",
		"-root",
		c.root,
		"-cache",
		c.cache,
		"-clean",
		"-debug",
	)
}
