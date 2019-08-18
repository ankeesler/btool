package integration

import (
	"github.com/stretchr/testify/require"
)

func version(c *config) {
	v := c.run(
		c.btool,
		"-version",
	)
	require.Regexp(c.t, `.*version \d\.\d\.\d.*`, v)
}
