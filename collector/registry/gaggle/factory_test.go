package gaggler_test

import (
	"testing"

	"github.com/ankeesler/btool/collector/registry/gaggler"
	"github.com/ankeesler/btool/collector/registry/gaggler/gagglerfakes"
	"github.com/ankeesler/btool/registry"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFactoryNext(t *testing.T) {
	// TODO: test caching

	fs := afero.NewMemMapFs()

	gaggleA := &registry.Gaggle{}
	gaggleB := &registry.Gaggle{}
	gaggleC := &registry.Gaggle{}
	r := &gagglerfakes.FakeRegistry{}
	r.IndexReturnsOnCall(
		0,
		&registry.Index{
			Files: []registry.IndexFile{
				registry.IndexFile{
					Path:   "a",
					SHA256: "a-sha",
				},
				registry.IndexFile{
					Path:   "b",
					SHA256: "b-sha",
				},
				registry.IndexFile{
					Path:   "c",
					SHA256: "c-sha",
				},
			},
		},
		nil,
	)
	r.GaggleReturnsOnCall(0, gaggleA, nil)
	r.GaggleReturnsOnCall(1, gaggleB, nil)
	r.GaggleReturnsOnCall(2, gaggleC, nil)

	cache := "/some/cache"
	f := gaggler.NewFactory(fs, r, cache)

	gagglers, err := drainGagglers(f)
	require.Nil(t, err)
	assert.Equal(t, 3, len(gagglers))
	assert.Equal(t, gaggleA, gagglers[0].Gaggle())
	assert.Equal(t, "/some/cache/a-sha", gagglers[0].Root())
	assert.Equal(t, gaggleB, gagglers[1].Gaggle())
	assert.Equal(t, "/some/cache/b-sha", gagglers[1].Root())
	assert.Equal(t, gaggleC, gagglers[2].Gaggle())
	assert.Equal(t, "/some/cache/c-sha", gagglers[2].Root())

	assert.Equal(t, 1, r.IndexCallCount())
	assert.Equal(t, 3, r.GaggleCallCount())
	assert.Equal(t, "a", r.GaggleArgsForCall(0))
	assert.Equal(t, "b", r.GaggleArgsForCall(1))
	assert.Equal(t, "c", r.GaggleArgsForCall(2))
}

func drainGagglers(f *gaggler.Factory) ([]*gaggler.Gaggler, error) {
	gagglers := make([]*gaggler.Gaggler, 0)
	for {
		gaggler, err := f.Next()
		if err != nil {
			return nil, errors.Wrap(err, "next")
		} else if gaggler == nil {
			return gagglers, nil
		} else {
			gagglers = append(gagglers, gaggler)
		}
	}
}
