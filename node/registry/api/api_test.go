package api_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/node/registry"
	"github.com/ankeesler/btool/node/registry/api"
	"github.com/ankeesler/btool/node/registry/api/apifakes"
	"github.com/ankeesler/btool/node/registry/testutil"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestAPI(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

	i := testutil.Index()
	iBytes, err := yaml.Marshal(i)
	require.Nil(t, err)
	nodes := testutil.FileANodes()
	nodesBytes, err := yaml.Marshal(nodes)
	require.Nil(t, err)

	r := &apifakes.FakeRegistry{}
	r.IndexReturnsOnCall(0, i, nil)
	r.NodesReturnsOnCall(0, nodes, nil)
	r.NodesReturnsOnCall(1, []*registry.Node{}, nil)

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

	// Nodes.
	nodesRsp, err := c.Get(s.URL + "/" + "file_a_btool.yml")
	require.Nil(t, err)
	defer nodesRsp.Body.Close()

	nodesData, err := ioutil.ReadAll(nodesRsp.Body)
	require.Nil(t, err)

	assert.Equal(t, 1, r.NodesCallCount())
	assert.Equal(t, "file_a_btool.yml", r.NodesArgsForCall(0))

	assert.Equal(t, http.StatusOK, nodesRsp.StatusCode)
	assert.Equal(t, nodesBytes, nodesData)

	rsp, err := c.Get(s.URL + "/" + "does not exist")
	require.Nil(t, err)
	defer rsp.Body.Close()

	assert.Equal(t, 2, r.NodesCallCount())
	assert.Equal(t, "does not exist", r.NodesArgsForCall(1))

	assert.Equal(t, http.StatusNotFound, rsp.StatusCode)
}
