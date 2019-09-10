package prodcon_test

import (
	"testing"

	"github.com/ankeesler/btool/collector/prodcon"
	"github.com/ankeesler/btool/collector/prodcon/prodconfakes"
	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPCCollect(t *testing.T) {
	p0 := &prodconfakes.FakeProducer{}
	p0n0 := node.New("p0n0")
	p0n1 := node.New("p0n1")
	p0.ProduceStub = func(s *prodcon.Store) error {
		s.Add(p0n0)
		s.Add(p0n1)
		return nil
	}
	p1 := &prodconfakes.FakeProducer{}
	p1n0 := node.New("p1n0")
	p1n1 := node.New("p1n1")
	p1.ProduceStub = func(s *prodcon.Store) error {
		s.Add(p1n0)
		s.Add(p1n1)
		return nil
	}
	producers := []prodcon.Producer{p0, p1}

	c0 := &prodconfakes.FakeConsumer{}
	c0n0 := node.New("c0n0")
	c0n1 := node.New("c0n1")
	c0.ConsumeStub = func(s *prodcon.Store, diff *prodcon.Diff) error {
		log.Debugf("c0.Consume %s", diff)
		assert.Equal(t, prodcon.DiffAdd, diff.Type)
		if diff.Name == "p0n0" {
			s.Add(c0n0)
		} else if diff.Name == "p0n1" {
			s.Add(c0n1)
		}
		return nil
	}
	c1 := &prodconfakes.FakeConsumer{}
	c1n0 := node.New("c1n0")
	c1n1 := node.New("c1n1")
	c1.ConsumeStub = func(s *prodcon.Store, diff *prodcon.Diff) error {
		log.Debugf("c1.Consume %s", diff)

		assert.Equal(t, prodcon.DiffAdd, diff.Type)
		if diff.Name == "p1n0" {
			s.Add(c1n0)
		} else if diff.Name == "p1n1" {
			s.Add(c1n1)
		}
		return nil
	}
	consumers := []prodcon.Consumer{c0, c1}

	acS := prodcon.NewStore()
	pc := prodcon.New(acS, producers, consumers)
	require.Nil(t, pc.Run())

	assert.Equal(t, 1, p0.ProduceCallCount())
	assert.Equal(t, 1, p1.ProduceCallCount())

	assert.Equal(t, 6, c0.ConsumeCallCount())
	assertDiff(t, c0, 0, "p0n0")
	assertDiff(t, c0, 1, "p0n1")
	assertDiff(t, c0, 2, "p1n0")
	assertDiff(t, c0, 3, "p1n1")
	assertDiff(t, c0, 4, "c1n0")
	assertDiff(t, c0, 5, "c1n1")
	assert.Equal(t, 6, c1.ConsumeCallCount())
	assertDiff(t, c1, 0, "p0n0")
	assertDiff(t, c1, 1, "p0n1")
	assertDiff(t, c1, 2, "p1n0")
	assertDiff(t, c1, 3, "p1n1")
	assertDiff(t, c1, 4, "c0n0")
	assertDiff(t, c1, 5, "c0n1")

	exS := prodcon.NewStore()
	exS.Add(p0n0)
	exS.Add(p0n1)
	exS.Add(p1n0)
	exS.Add(p1n1)
	exS.Add(c0n0)
	exS.Add(c0n1)
	exS.Add(c1n0)
	exS.Add(c1n1)
	assert.Nil(t, deep.Equal(exS, acS))
}

func assertDiff(
	t *testing.T,
	c *prodconfakes.FakeConsumer,
	call int,
	exName string,
) {
	exDiff := &prodcon.Diff{
		Type: prodcon.DiffAdd,
		Name: exName,
	}
	_, acDiff := c.ConsumeArgsForCall(call)
	assert.Equal(t, exDiff, acDiff)
}
