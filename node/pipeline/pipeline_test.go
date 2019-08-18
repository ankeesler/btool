package pipeline_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/ankeesler/btool/node/pipeline"
	"github.com/ankeesler/btool/node/pipeline/pipelinefakes"
)

func TestPipeline(t *testing.T) {
	goodHandlerA := &pipelinefakes.FakeHandler{}
	goodHandlerA.HandleStub = func(ctx *pipeline.Ctx) error {
		ctx.KV["key"] = "value"
		return nil
	}
	goodHandlerB := &pipelinefakes.FakeHandler{}
	goodHandlerB.HandleStub = func(ctx *pipeline.Ctx) error {
		if ctx.KV["key"] != "value" {
			t.Error("expected kv entry")
		}
		return nil
	}
	badHandler := &pipelinefakes.FakeHandler{}
	badHandler.HandleStub = func(ctx *pipeline.Ctx) error {
		return errors.New("some error")
	}

	ctx := pipeline.NewCtx()

	// Happy
	goodP := pipeline.New(ctx, goodHandlerA, goodHandlerB)
	if err := goodP.Run(); err != nil {
		t.Error(err)
	}

	// Sad.
	badP := pipeline.New(ctx, goodHandlerA, badHandler, goodHandlerB)
	if err := badP.Run(); err == nil {
		t.Error("expected failure")
	}

	data := []struct {
		ex              int
		callCountFunc   func() int
		argsForCallFunc func(int) *pipeline.Ctx
	}{
		{
			ex:              2,
			callCountFunc:   goodHandlerA.HandleCallCount,
			argsForCallFunc: goodHandlerA.HandleArgsForCall,
		},
		{
			ex:              1,
			callCountFunc:   goodHandlerB.HandleCallCount,
			argsForCallFunc: goodHandlerB.HandleArgsForCall,
		},
		{
			ex:              1,
			callCountFunc:   badHandler.HandleCallCount,
			argsForCallFunc: badHandler.HandleArgsForCall,
		},
	}
	for _, datum := range data {
		if ac := datum.callCountFunc(); datum.ex != ac {
			t.Errorf("%s: %d != %d", reflect.TypeOf(datum.callCountFunc).Name(), datum.ex, ac)
		}

		ex := ctx
		for i := 0; i < datum.callCountFunc(); i++ {
			ac := datum.argsForCallFunc(i)
			if ex != ac {
				t.Error(ex, "!=", ac)
			}
		}
	}
}
