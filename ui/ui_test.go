package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortenName(t *testing.T) {
	data := []struct {
		in, out string
	}{
		{
			// short enough
			in:  "abc123/a.cc",
			out: "abc123/a.cc",
		},
		{
			// long
			in:  "0123456789abcdefghij/0123456789.txt",
			out: "012345678.../0123456789.txt",
		},
	}
	for _, datum := range data {
		assert.Equal(t, datum.out, shortenName(datum.in))
	}
}
