package registry_test

import (
	"path/filepath"
	"testing"

	"github.com/ankeesler/btool/collector"
	"github.com/ankeesler/btool/collector/registry"
	"github.com/ankeesler/btool/collector/registry/registryfakes"
	"github.com/ankeesler/btool/collector/testutil"
	registrypkg "github.com/ankeesler/btool/registry"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistryCollect(t *testing.T) {
	// TODO: test caching

	fs := afero.NewMemMapFs()

	gaggleA := &registrypkg.Gaggle{}
	gaggleB := &registrypkg.Gaggle{}
	gaggleC := &registrypkg.Gaggle{}
	c := &registryfakes.FakeClient{}
	c.IndexReturnsOnCall(
		0,
		&registrypkg.Index{
			Files: []registrypkg.IndexFile{
				registrypkg.IndexFile{
					Path:   "a",
					SHA256: "a-sha",
				},
				registrypkg.IndexFile{
					Path:   "b",
					SHA256: "b-sha",
				},
				registrypkg.IndexFile{
					Path:   "c",
					SHA256: "c-sha",
				},
			},
		},
		nil,
	)
	c.GaggleReturnsOnCall(0, gaggleA, nil)
	c.GaggleReturnsOnCall(1, gaggleB, nil)
	c.GaggleReturnsOnCall(2, gaggleC, nil)

	cache := "/some/cache"
	gc := &registryfakes.FakeGaggleCollector{}
	r := registry.New(fs, c, cache, gc)

	s := testutil.FakeStore()
	require.Nil(t, r.Produce(s))

	require.Equal(t, 1, c.IndexCallCount())
	require.Equal(t, 3, c.GaggleCallCount())
	assert.Equal(t, "a", c.GaggleArgsForCall(0))
	assert.Equal(t, "b", c.GaggleArgsForCall(1))
	assert.Equal(t, "c", c.GaggleArgsForCall(2))

	var acS collector.Store
	var acGaggle *registrypkg.Gaggle
	var acRoot string
	require.Equal(t, 3, gc.CollectCallCount())
	acS, acGaggle, acRoot = gc.CollectArgsForCall(0)
	assert.Equal(t, s, acS)
	assert.Equal(t, gaggleA, acGaggle)
	assert.Equal(t, filepath.Join(cache, "a-sha"), acRoot)
	acS, acGaggle, acRoot = gc.CollectArgsForCall(1)
	assert.Equal(t, s, acS)
	assert.Equal(t, gaggleB, acGaggle)
	assert.Equal(t, filepath.Join(cache, "b-sha"), acRoot)
	acS, acGaggle, acRoot = gc.CollectArgsForCall(2)
	assert.Equal(t, s, acS)
	assert.Equal(t, gaggleC, acGaggle)
	assert.Equal(t, filepath.Join(cache, "c-sha"), acRoot)
}
