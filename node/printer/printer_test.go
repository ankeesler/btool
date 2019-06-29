package printer_test

import (
	"bytes"
	"testing"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/printer"
	"github.com/ankeesler/btool/node/testutil"
	"github.com/go-test/deep"
)

func TestHandle(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	p := printer.New(buf)

	nodes, err := p.Handle(testutil.BasicNodesC)
	if err != nil {
		t.Error(err)
	}

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

	if diff := deep.Equal(nodes, []*node.Node(testutil.BasicNodesC)); diff != nil {
		t.Error(diff)
	}
}
