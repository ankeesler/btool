package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortenName(t *testing.T) {
	data := []struct {
		cache, in, out string
	}{
		{
			// short enough
			cache: "nah",
			in:    "abc123/a.cc",
			out:   "abc123/a.cc",
		},
		{
			// short enough with cache
			cache: "abc123",
			in:    "abc123/a.cc",
			out:   "$CACHE/a.cc",
		},
		{
			// long
			cache: "nah",
			in:    "0123456789abcdefghij/0123456789.txt",
			out:   "0123456789.../0123456789.txt",
		},
		{
			// long with cache
			cache: "a",
			in:    "a/0123456789abcdefghij/0123456789.txt",
			out:   "$CACHE/012.../0123456789.txt",
		},
	}
	for _, datum := range data {
		assert.Equal(t, datum.out, shortenName(datum.in, datum.cache))
	}
}
