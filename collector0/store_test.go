package collector

import (
	"testing"

	"github.com/ankeesler/btool/node"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	s := newStore()
	n0 := node.New("0")
	n1 := node.New("1")
	n2 := node.New("2")
	s.Set(n0)
	s.Set(n1)
	assert.Equal(t, n0, s.Get(n0.Name))
	assert.Equal(t, n1, s.Get(n1.Name))
	assert.Nil(t, s.Get(n2.Name))

	acNodes := make(map[string]*node.Node)
	s.ForEach(func(n *node.Node) {
		acNodes[n.Name] = n
	})
	exNodes := map[string]*node.Node{
		"0": n0,
		"1": n1,
	}
	assert.Equal(t, exNodes, acNodes)
}
