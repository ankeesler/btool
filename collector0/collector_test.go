package collector_test

import (
	"testing"

	collector "github.com/ankeesler/btool/collector0"
	collectorfakes "github.com/ankeesler/btool/collector0/collector0fakes"
	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectorCollect(t *testing.T) {
	acN := node.New("n")

	p0 := &collectorfakes.FakeProducer{}
	p0n0 := node.New("p0n0")
	p0n1 := node.New("p0n1")
	p0.ProduceStub = func(s collector.Store) error {
		s.Set(p0n0)
		s.Set(p0n1)
		acN.Dependency(p0n0, p0n1)
		return nil
	}
	p1 := &collectorfakes.FakeProducer{}
	p1n0 := node.New("p1n0")
	p1n1 := node.New("p1n1")
	p1.ProduceStub = func(s collector.Store) error {
		s.Set(p1n0)
		s.Set(p1n1)
		acN.Dependency(p1n0, p1n1)
		return nil
	}
	producers := []collector.Producer{p0, p1}

	c0 := &collectorfakes.FakeConsumer{}
	c0n0 := node.New("c0n0")
	c0n1 := node.New("c0n1")
	c0.ConsumeStub = func(s collector.Store, n *node.Node) error {
		log.Debugf("c0.Consume %s", n)
		if n.Name == "p0n0" {
			s.Set(c0n0)
			acN.Dependency(c0n0)
		} else if n.Name == "p0n1" {
			s.Set(c0n1)
			acN.Dependency(c0n1)
		}
		return nil
	}
	c1 := &collectorfakes.FakeConsumer{}
	c1n0 := node.New("c1n0")
	c1n1 := node.New("c1n1")
	c1.ConsumeStub = func(s collector.Store, n *node.Node) error {
		log.Debugf("c1.Consume %s", n)
		if n.Name == "p1n0" {
			s.Set(c1n0)
			acN.Dependency(c1n0)
		} else if n.Name == "p1n1" {
			s.Set(c1n1)
			acN.Dependency(c1n1)
		}
		return nil
	}
	consumers := []collector.Consumer{c0, c1}

	c := collector.New(producers, consumers)
	require.Nil(t, c.Collect(acN))

	assert.Equal(t, 1, p0.ProduceCallCount())
	assert.Equal(t, 1, p1.ProduceCallCount())

	assert.Equal(t, 6, c0.ConsumeCallCount())
	assertConsumeArgs(t, c0, 0, p0n0)
	assertConsumeArgs(t, c0, 1, p0n1)
	assertConsumeArgs(t, c0, 2, p1n0)
	assertConsumeArgs(t, c0, 3, p1n1)
	assertConsumeArgs(t, c0, 4, c1n0)
	assertConsumeArgs(t, c0, 5, c1n1)
	assert.Equal(t, 6, c1.ConsumeCallCount())
	assertConsumeArgs(t, c1, 0, p0n0)
	assertConsumeArgs(t, c1, 1, p0n1)
	assertConsumeArgs(t, c1, 2, p1n0)
	assertConsumeArgs(t, c1, 3, p1n1)
	assertConsumeArgs(t, c1, 4, c0n0)
	assertConsumeArgs(t, c1, 5, c0n1)

	exN := node.New("n").Dependency(
		p0n0,
		p0n1,
		p1n0,
		p1n1,
		c0n0,
		c0n1,
		c1n0,
		c1n1,
	)
	assert.Equal(t, exN, acN)
}

func assertConsumeArgs(
	t *testing.T,
	c *collectorfakes.FakeConsumer,
	call int,
	exN *node.Node,
) {
	_, acN := c.ConsumeArgsForCall(call)
	assert.Equal(t, exN, acN)
}
