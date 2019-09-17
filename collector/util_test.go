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
				d := node.New("d").Label("l", []string{"l-d"})
				c := node.New("c").Dependency(d).Label("l", []string{"l-c"})
				b := node.New("b").Label("l", []string{"l-b"})
				a := node.New("a").Dependency(b, c).Label("l", []string{"l-a"})
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

func TestBoolLabel(t *testing.T) {
	data := []struct {
		name string

		n     *node.Node
		label string

		exBool    bool
		exSuccess bool
	}{
		{
			name:      "Unknown",
			n:         node.New("a"),
			label:     "l",
			exBool:    false,
			exSuccess: true,
		},
		{
			name:      "WrongType",
			n:         node.New("a").Label("l", "hey"),
			label:     "l",
			exBool:    false,
			exSuccess: false,
		},
		{
			name:      "Success",
			n:         node.New("a").Label("l", true),
			label:     "l",
			exBool:    true,
			exSuccess: true,
		},
	}
	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			acBool, err := collector.BoolLabel(datum.n, datum.label)
			acSuccess := err == nil
			assert.Equal(t, datum.exBool, acBool)
			assert.Equal(t, datum.exSuccess, acSuccess)
		})
	}
}

func TestAppendLabel(t *testing.T) {
	data := []struct {
		name string

		n     *node.Node
		label string
		apend string

		exLabels  []string
		exSuccess bool
	}{
		{
			name:      "NewLabel",
			n:         node.New("a"),
			label:     "l",
			apend:     "hey",
			exLabels:  []string{"hey"},
			exSuccess: true,
		},
		{
			name:      "ExistingLabel",
			n:         node.New("a").Label("l", []string{"hey"}),
			label:     "l",
			apend:     "ho",
			exLabels:  []string{"hey", "ho"},
			exSuccess: true,
		},
	}
	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			err := collector.AppendLabel(datum.n, datum.label, datum.apend)
			acSuccess := err == nil
			assert.Equal(t, datum.exSuccess, acSuccess)

			assert.Equal(t, datum.exLabels, datum.n.Labels[datum.label])
		})
	}
}
