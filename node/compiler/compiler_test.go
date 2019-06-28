package compiler_test

import (
	"testing"

	"github.com/ankeesler/btool/node/compiler"
	"github.com/ankeesler/btool/node/compiler/compilerfakes"
	"github.com/ankeesler/btool/node/testutil"
	"github.com/go-test/deep"
)

func TestHandle(t *testing.T) {
	c := &compilerfakes.FakeC{}
	compiler := compiler.New(c, "/", "/cache")

	nodes := testutil.BasicNodes
	exNodes := testutil.AddObjects(nodes)

	acNodes, err := compiler.Handle(nodes)
	if err != nil {
		t.Fatal(err)
	}

	if diff := deep.Equal(exNodes, acNodes); diff != nil {
		t.Error(diff)
	}
}
