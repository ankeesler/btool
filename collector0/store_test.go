package collector_test

import (
	"testing"

	collector "github.com/ankeesler/btool/collector0"
	"github.com/ankeesler/btool/node"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	s := collector.NewStore()
	n0 := node.New("0")
	n1 := node.New("1")
	n2 := node.New("2")
	s.Add(n0, n1)
	assert.Equal(t, n0, s.Find(n0.Name))
	assert.Equal(t, n1, s.Find(n1.Name))
	assert.Nil(t, s.Find(n2.Name))
}
