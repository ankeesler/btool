package collector_test

import (
	"testing"

	collector "github.com/ankeesler/btool/collector0"
	collectorfakes "github.com/ankeesler/btool/collector0/collector0fakes"
	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectorCollect(t *testing.T) {
	p0 := &collectorfakes.FakeProducer{}
	p0n0 := node.New("p0n0")
	p0n1 := node.New("p0n1")
	p0.ProduceStub = func(s *collector.Store) error {
		s.Add(p0n0)
		s.Add(p0n1)
		return nil
	}
	p1 := &collectorfakes.FakeProducer{}
	p1n0 := node.New("p1n0")
	p1n1 := node.New("p1n1")
	p1.ProduceStub = func(s *collector.Store) error {
		s.Add(p1n0)
		s.Add(p1n1)
		return nil
	}
	producers := []collector.Producer{p0, p1}

	c0 := &collectorfakes.FakeConsumer{}
	c0n0 := node.New("c0n0")
	c0n1 := node.New("c0n1")
	c0.ConsumeStub = func(s *collector.Store, diff *collector.Diff) error {
		log.Debugf("c0.Consume %s", diff)
		assert.Equal(t, collector.DiffAdd, diff.Type)
		if diff.Name == "p0n0" {
			s.Add(c0n0)
		} else if diff.Name == "p0n1" {
			s.Add(c0n1)
		}
		return nil
	}
	c1 := &collectorfakes.FakeConsumer{}
	c1n0 := node.New("c1n0")
	c1n1 := node.New("c1n1")
	c1.ConsumeStub = func(s *collector.Store, diff *collector.Diff) error {
		log.Debugf("c1.Consume %s", diff)

		assert.Equal(t, collector.DiffAdd, diff.Type)
		if diff.Name == "p1n0" {
			s.Add(c1n0)
		} else if diff.Name == "p1n1" {
			s.Add(c1n1)
		}
		return nil
	}
	consumers := []collector.Consumer{c0, c1}

	acS := collector.NewStore()
	c := collector.New(acS, producers, consumers)
	require.Nil(t, c.Collect(nil))

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

	exS := collector.NewStore()
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
	c *collectorfakes.FakeConsumer,
	call int,
	exName string,
) {
	exDiff := &collector.Diff{
		Type: collector.DiffAdd,
		Name: exName,
	}
	_, acDiff := c.ConsumeArgsForCall(call)
	assert.Equal(t, exDiff, acDiff)
}
