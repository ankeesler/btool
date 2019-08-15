package handlers_test

import (
	"path/filepath"
	"testing"

	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/nodefakes"
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
	s.RegistryDirReturns("/some-registry-dir")
	s.ProjectDirReturns("/some-project-dir")

	compileCR := &nodefakes.FakeResolver{}
	compileCCR := &nodefakes.FakeResolver{}
	unzipR := &nodefakes.FakeResolver{}
	downloadR := &nodefakes.FakeResolver{}

	rf := &handlersfakes.FakeResolverFactory{}
	rf.NewCompileCReturnsOnCall(0, compileCR)
	rf.NewCompileCCReturnsOnCall(0, compileCCR)
	rf.NewUnzipReturnsOnCall(0, unzipR)
	rf.NewDownloadReturnsOnCall(0, downloadR)

	r := &handlersfakes.FakeRegistry{}
	index := testutil.Index()
	registryFileAGaggle := testutil.FileAGaggle()
	registryFileAGaggle.Nodes[0].Resolver.Name = "compileC"
	registryFileAGaggle.Nodes[1].Resolver.Name = "compileCC"
	registryFileBGaggle := testutil.FileBGaggle()
	registryFileBGaggle.Nodes[0].Resolver.Name = "unzip"
	registryFileBGaggle.Nodes[1].Resolver.Name = "download"
	registryFileBGaggle.Nodes[1].Resolver.Config["url"] = "some url"
	registryFileBGaggle.Nodes[1].Resolver.Config["sha256"] = "some sha"
	r.IndexReturns(index, nil)
	r.GaggleReturnsOnCall(0, registryFileAGaggle, nil)

	data, err := yaml.Marshal(registryFileBGaggle)
	assert.Nil(t, err)
	assert.Nil(t, afero.WriteFile(
		fs,
		filepath.Join(
			"/some-registry-dir",
			index.Files[1].SHA256+".yml",
		),
		data,
		0644,
	))

	tunaN := node.New("/some-project-dir/tuna")
	tunaN.Resolver = compileCR
	fishN := node.New("/some-project-dir/fish").Dependency(tunaN)
	fishN.Resolver = compileCCR
	marlinN := node.New("/some-project-dir/marlin")
	marlinN.Resolver = unzipR
	baconN := node.New("/some-project-dir/bacon").Dependency(marlinN)
	baconN.Resolver = downloadR

	h := handlers.NewRegistry(fs, s, rf, r)
	ctx := pipeline.NewCtx()
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

	assert.Equal(t, 1, rf.NewCompileCCallCount())
	assert.Equal(t, []string{"/some-project-dir"}, rf.NewCompileCArgsForCall(0))
	assert.Equal(t, 1, rf.NewCompileCCCallCount())
	assert.Equal(t, []string{"/some-project-dir"}, rf.NewCompileCCArgsForCall(0))
	assert.Equal(t, 1, rf.NewUnzipCallCount())
	assert.Equal(t, "/some-project-dir", rf.NewUnzipArgsForCall(0))
	assert.Equal(t, 1, rf.NewDownloadCallCount())
	url, sha256 := rf.NewDownloadArgsForCall(0)
	assert.Equal(t, "some url", url)
	assert.Equal(t, "some sha", sha256)

	assert.Equal(t, 1, r.IndexCallCount())
	assert.Equal(t, 1, r.GaggleCallCount())
	assert.Equal(t, index.Files[0].Path, r.GaggleArgsForCall(0))

	exNodes := []*node.Node{tunaN, fishN, marlinN, baconN}
	assert.Nil(t, deep.Equal(exNodes, ctx.Nodes))
}
