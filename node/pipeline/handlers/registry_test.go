package handlers_test

//func TestRegistry(t *testing.T) {
//	logrus.SetLevel(logrus.DebugLevel)
//	logrus.SetFormatter(formatter.New())
//
//	fs := afero.NewMemMapFs()
//
//	index := testutil.Index()
//	registryFileANodes := testutil.FileANodes()
//	registryFileBNodes := testutil.FileBNodes()
//	tunaN := node.New("tuna")
//	fishN := node.New("fish")
//	marlinN := node.New("marlin")
//	baconN := node.New("bacon")
//
//	cache := "/cache"
//	data, err := yaml.Marshal(registryFileBNodes)
//	assert.Nil(t, err)
//	assert.Nil(t, afero.WriteFile(
//		fs,
//		filepath.Join(
//			cache,
//			"some/path/to/file_b",
//			"download",
//			index.Files[1].SHA256,
//		),
//		data,
//		0644,
//	))
//
//	r := &handlersfakes.FakeRegistry{}
//	r.IndexReturnsOnCall(0, index, nil)
//	r.NodesReturnsOnCall(0, registryFileANodes, nil)
//	d := &handlersfakes.FakeDecoder{}
//	d.DecodeReturnsOnCall(0, tunaN, nil)
//	d.DecodeReturnsOnCall(1, fishN, nil)
//	d.DecodeReturnsOnCall(2, marlinN, nil)
//	d.DecodeReturnsOnCall(3, baconN, nil)
//
//	h := handlers.NewRegistry(fs, r, d)
//	ctx := pipeline.NewCtxBuilder().Cache(cache).Build()
//	assert.Nil(t, h.Handle(ctx))
//
//	assert.Equal(t, 1, r.IndexCallCount())
//
//	assert.Equal(t, 1, r.NodesCallCount())
//	for i := 0; i < r.NodesCallCount(); i++ {
//		assert.Equal(t, index.Files[i].Path, r.NodesArgsForCall(i))
//	}
//
//	assert.Equal(t, 4, d.DecodeCallCount())
//	assert.Equal(t, registryFileANodes[0], d.DecodeArgsForCall(0))
//	assert.Equal(t, registryFileANodes[1], d.DecodeArgsForCall(1))
//	assert.Equal(t, registryFileBNodes[0], d.DecodeArgsForCall(2))
//	assert.Equal(t, registryFileBNodes[1], d.DecodeArgsForCall(3))
//
//	exNodes := []*node.Node{tunaN, fishN, marlinN, baconN}
//	assert.Equal(t, exNodes, ctx.Nodes)
//}
