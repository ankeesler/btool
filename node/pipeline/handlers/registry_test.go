package handlers_test

import (
	"path/filepath"
	"testing"

	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/pipeline"
	"github.com/ankeesler/btool/node/pipeline/handlers"
	"github.com/ankeesler/btool/node/pipeline/handlers/handlersfakes"
	"github.com/ankeesler/btool/node/registry/testutil"
	"github.com/go-test/deep"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestRegistry(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

	fs := afero.NewMemMapFs()

	s := &handlersfakes.FakeStore{}
	s.RegistryDirReturns("/some-registry-path")
	s.ProjectDirReturns("/some-project-path")

	r := &handlersfakes.FakeRegistry{}
	index := testutil.Index()
	registryFileAGaggle := testutil.FileAGaggle()
	registryFileBGaggle := testutil.FileBGaggle()
	r.IndexReturns(index, nil)
	r.GaggleReturnsOnCall(0, registryFileAGaggle, nil)

	data, err := yaml.Marshal(registryFileBGaggle)
	assert.Nil(t, err)
	assert.Nil(t, afero.WriteFile(
		fs,
		filepath.Join(
			"/some-registry-path",
			index.Files[1].SHA256+".yml",
		),
		data,
		0644,
	))

	tunaN := node.New("/some-project-path/tuna")
	fishN := node.New("/some-project-path/fish").Dependency(tunaN)
	marlinN := node.New("/some-project-path/marlin")
	baconN := node.New("/some-project-path/bacon").Dependency(marlinN)

	h := handlers.NewRegistry(fs, s, r)
	ctx := pipeline.NewCtxBuilder().Build()
	assert.Nil(t, h.Handle(ctx))

	assert.Equal(t, 1, s.RegistryDirCallCount())
	assert.Equal(t, index.Name, s.RegistryDirArgsForCall(0))
	assert.Equal(t, 2, s.ProjectDirCallCount())
	assert.Equal(
		t,
		registryFileAGaggle.Metadata["project"],
		s.ProjectDirArgsForCall(0),
	)
	assert.Equal(
		t,
		registryFileBGaggle.Metadata["project"],
		s.ProjectDirArgsForCall(1),
	)

	assert.Equal(t, 1, r.IndexCallCount())
	assert.Equal(t, 1, r.GaggleCallCount())
	assert.Equal(t, index.Files[0].Path, r.GaggleArgsForCall(0))

	exNodes := []*node.Node{tunaN, fishN, marlinN, baconN}
	assert.Nil(t, deep.Equal(exNodes, ctx.Nodes))
}
