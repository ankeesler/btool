package pipeline_test

import (
	"errors"
	"testing"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/pipeline"
	"github.com/ankeesler/btool/node/pipeline/pipelinefakes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPipeline(t *testing.T) {
	nodeA := node.New("a")
	nodeB := node.New("b").Dependency(nodeA)
	goodHandlerA := &pipelinefakes.FakeHandler{}
	goodHandlerA.HandleStub = func(ctx pipeline.Ctx) error {
		ctx.Add(nodeA)
		return nil
	}
	goodHandlerB := &pipelinefakes.FakeHandler{}
	goodHandlerB.HandleStub = func(ctx pipeline.Ctx) error {
		assert.NotNil(t, ctx.Find(nodeA.Name))
		ctx.Add(nodeB)
		return nil
	}
	badHandler := &pipelinefakes.FakeHandler{}
	badHandler.HandleStub = func(ctx pipeline.Ctx) error {
		return errors.New("some error")
	}

	// Happy
	cb := &pipelinefakes.FakeCallback{}
	mh := pipeline.NewMultiHandler().Add(goodHandlerA, goodHandlerB)
	goodP := pipeline.New(mh, cb)
	acNodes, err := goodP.Run()
	require.Nil(t, err)
	assert.Equal(
		t,
		[]*node.Node{
			nodeA, nodeB,
		},
		acNodes,
	)
	assert.Equal(t, 2, cb.OnAddCallCount())
	assert.Equal(t, nodeA, cb.OnAddArgsForCall(0))
	assert.Equal(t, nodeB, cb.OnAddArgsForCall(1))

	// Sad.
	cb = &pipelinefakes.FakeCallback{}
	mh = pipeline.NewMultiHandler().Add(goodHandlerA, badHandler)
	badP := pipeline.New(mh, cb)
	_, err = badP.Run()
	assert.NotNil(t, err)
	assert.Equal(t, 1, cb.OnAddCallCount())
	assert.Equal(t, nodeA, cb.OnAddArgsForCall(0))

	data := []struct {
		ex            int
		callCountFunc func() int
	}{
		{
			ex:            2,
			callCountFunc: goodHandlerA.HandleCallCount,
		},
		{
			ex:            1,
			callCountFunc: goodHandlerB.HandleCallCount,
		},
		{
			ex:            1,
			callCountFunc: badHandler.HandleCallCount,
		},
	}
	for _, datum := range data {
		assert.Equal(t, datum.ex, datum.callCountFunc())
	}
}
