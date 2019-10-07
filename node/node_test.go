package node_test

import (
	"testing"

	"github.com/ankeesler/btool/node"
	"github.com/stretchr/testify/assert"
)

func TestNodeDependency(t *testing.T) {
	a := node.New("a")
	b0 := node.New("b")
	b1 := node.New("b")
	a.Dependency(b0)
	a.Dependency(b1)
	assert.Equal(t, 1, len(a.Dependencies))
	assert.Equal(t, b1, a.Dependencies[0])
}
