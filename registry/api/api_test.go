package api_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankeesler/btool/registry"
	"github.com/ankeesler/btool/registry/api"
	"github.com/ankeesler/btool/registry/api/apifakes"
	"github.com/ankeesler/btool/registry/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestAPI(t *testing.T) {
	i := testutil.Index()
	iBytes, err := yaml.Marshal(i)
	require.Nil(t, err)
	gaggle := testutil.FileAGaggle()
	gaggleBytes, err := yaml.Marshal(gaggle)
	require.Nil(t, err)

	r := &apifakes.FakeRegistry{}
	r.IndexReturnsOnCall(0, i, nil)
	r.GaggleStub = func(name string) (*registry.Gaggle, error) {
		if r.GaggleCallCount() == 1 {
			return gaggle, nil
		} else {
			return nil, nil
		}
	}

	a := api.New(r)
	s := httptest.NewServer(a)
	defer s.Close()

	c := s.Client()

	// Index.
	iRsp, err := c.Get(s.URL)
	require.Nil(t, err)
	defer iRsp.Body.Close()

	iData, err := ioutil.ReadAll(iRsp.Body)
	require.Nil(t, err)

	assert.Equal(t, 1, r.IndexCallCount())

	assert.Equal(t, http.StatusOK, iRsp.StatusCode)
	assert.Equal(t, iBytes, iData)

	// Gaggle.
	gaggleRsp, err := c.Get(s.URL + "/" + "file_a_btool.yml")
	require.Nil(t, err)
	defer gaggleRsp.Body.Close()

	gaggleData, err := ioutil.ReadAll(gaggleRsp.Body)
	require.Nil(t, err)

	assert.Equal(t, 1, r.GaggleCallCount())
	assert.Equal(t, "file_a_btool.yml", r.GaggleArgsForCall(0))

	assert.Equal(t, http.StatusOK, gaggleRsp.StatusCode)
	assert.Equal(t, gaggleBytes, gaggleData)

	rsp, err := c.Get(s.URL + "/" + "does not exist")
	require.Nil(t, err)
	defer rsp.Body.Close()

	assert.Equal(t, 2, r.GaggleCallCount())
	assert.Equal(t, "does not exist", r.GaggleArgsForCall(1))

	assert.Equal(t, http.StatusNotFound, rsp.StatusCode)
}
