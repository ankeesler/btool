package collector_test

import (
	"testing"

	"github.com/ankeesler/btool/collector"
	"github.com/ankeesler/btool/node"
	"github.com/stretchr/testify/assert"
)

type S struct {
	A int      `mapstructure:"a"`
	B []string `mapstructure:"b"`
	C bool     `mapstructure:"c"`
}

func TestLabels(t *testing.T) {
	t.Run("ToFrom", func(t *testing.T) {
		st := S{
			A: 5,
			B: []string{"one", "two"},
			C: true,
		}
		n := node.New("n")
		collector.MustToLabels(n, &st)
		assert.Equal(t, n.Labels, map[string]interface{}{
			"a": 5,
			"b": []string{"one", "two"},
			"c": true,
		})

		sf := S{}
		collector.MustFromLabels(n, &sf)
		assert.Equal(t, st, sf)
	})

	t.Run("FromToSubset", func(t *testing.T) {
		nf := node.New("n").Label("a", 10).Label("b", []string{"hey", "ho"})
		s := S{}
		collector.MustFromLabels(nf, &s)
		assert.Equal(t, s, S{A: 10, B: []string{"hey", "ho"}, C: false})

		nt := node.New("n")
		collector.MustToLabels(nt, &s)
		assert.Equal(t, nf.Label("c", false), nt)
	})

	t.Run("FromToSuperset", func(t *testing.T) {
		acN := node.New("n").Label(
			"a", 15,
		).Label(
			"b", []string{"hey", "ho"},
		).Label(
			"c", true,
		).Label(
			"d", "whatever",
		)

		s := S{}
		collector.MustFromLabels(acN, &s)
		assert.Equal(t, s, S{A: 15, B: []string{"hey", "ho"}, C: true})
		s.B = append(s.B, "let's go")

		collector.MustToLabels(acN, &s)
		exN := node.New("n").Label(
			"a", 15,
		).Label(
			"b", []string{"hey", "ho", "let's go"},
		).Label(
			"c", true,
		).Label(
			"d", "whatever",
		)
		assert.Equal(t, exN, acN)
	})
}
