package handlers_test

import (
	"bytes"
	"testing"

	"github.com/ankeesler/btool/node/testutil"
	"github.com/ankeesler/btool/pipeline"
	"github.com/ankeesler/btool/pipeline/handlers"
	"github.com/go-test/deep"
)

func TestPrint(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	h := handlers.NewPrint(buf)

	ctx := pipeline.NewCtxBuilder().Nodes(
		testutil.BasicNodesC.Copy(),
	).Root(
		"/",
	).Build()
	h.Handle(ctx)
	if ctx.Err != nil {
		t.Error(ctx.Err)
	}

	ex := `*** Nodes ***
dep-0/dep-0.c
> dep-0/dep-0.h
dep-0/dep-0.h
dep-1/dep-1.c
> dep-1/dep-1.h
> dep-0/dep-0.h
dep-1/dep-1.h
> dep-0/dep-0.h
main.c
> dep-1/dep-1.h
> dep-0/dep-0.h
*** Err ***
<nil>
*** KV ***
map[pipeline.root:/]
`
	ac := buf.String()
	if ex != ac {
		t.Error(ex, "!=", ac)
	}

	if diff := deep.Equal(ctx.Nodes, testutil.BasicNodesC.Cast()); diff != nil {
		t.Error(diff)
	}
}
