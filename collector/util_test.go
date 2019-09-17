package collector_test

import (
	"testing"

	"github.com/ankeesler/btool/collector"
	"github.com/ankeesler/btool/node"
	"github.com/stretchr/testify/assert"
)

func TestCollectLabels(t *testing.T) {
	data := []struct {
		name string

		nFunc func() *node.Node
		label string

		exLabels  []string
		exSuccess bool
	}{
		{
			name: "Simple",

			nFunc: func() *node.Node {
				// a -> b, c
				// b
				// c -> d
				// d
				d := node.New("d").Label("l", "l-d")
				c := node.New("c").Dependency(d).Label("l", "l-c")
				b := node.New("b").Label("l", "l-b")
				a := node.New("a").Dependency(b, c).Label("l", "l-a")
				return a
			},
			label: "l",

			exLabels: []string{
				"l-b",
				"l-d",
				"l-c",
				"l-a",
			},
			exSuccess: true,
		},
	}
	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			n := datum.nFunc()
			acLabels, err := collector.CollectLabels(n, datum.label)
			acSuccess := err == nil
			assert.Equal(t, datum.exLabels, acLabels)
			assert.Equal(t, datum.exSuccess, acSuccess)
		})
	}
}
