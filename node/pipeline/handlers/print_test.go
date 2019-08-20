package handlers_test

import (
	"bytes"
	"testing"

	"github.com/ankeesler/btool/node/pipeline/handlers"
	pipelinetestutil "github.com/ankeesler/btool/node/pipeline/testutil"
	"github.com/ankeesler/btool/node/testutil"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
)

func TestPrint(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	h := handlers.NewPrint(buf)

	ctx := pipelinetestutil.NewCtx()
	ctx.Nodes = testutil.BasicNodesC.Copy()
	assert.Nil(t, h.Handle(ctx))

	ex := `dep-0/dep-0.c
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
`
	ac := buf.String()
	if ex != ac {
		t.Error(ex, "!=", ac)
	}

	if diff := deep.Equal(ctx.All(), testutil.BasicNodesC.Cast()); diff != nil {
		t.Error(diff)
	}
}
