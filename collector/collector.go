// Package collector provides a producer/consumer framework that can act as a
// node.Node graph builder.
//
// Producer's add node.Node's to a Store. Consumer's are notified of node.Node's
// being Set() to the Store. A Store is just a place where node.Node's are kept.
package collector

import (
	"fmt"

	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
)

// These constants are node.Node Label's that are used throughout this framework.
const (
	LabelLocal = "io.btool.local"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Store

// Store is a place where node.Node's are kept.
//
// It has a simple idempotent interface. Set()'ing a node.Node for the first time
// will create it, setting a node.Node for the second time will update it.
type Store interface {
	Get(string) *node.Node
	Set(*node.Node)
	ForEach(func(*node.Node))
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Producer

// Producer is a type that adds node.Node's to a Store.
type Producer interface {
	Produce(Store) error
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Consumer

// Consumer is a type that reacts to Producer's adding node.Node's to a Store.
// The Consumer is provided the node.Node that was Set() on the Store.
type Consumer interface {
	Consume(Store, *node.Node) error
}

// Collector is an object that will run Producer's and Consumer's in order to
// build a node.Node graph.
//
// This particular procedure is very synchronus. It will:
//   1. run each Producer to completion
//   2. run all Consumer's on every Store.Set() from a Producer
//   3. run all Consumer's on every Store.Set() from a Consumer that isn't their own
//   4. repeat 4 until there are no more new Store.Set() calls
type Collector struct {
	s         *store
	producers []Producer
	consumers []Consumer
}

// New creates a new Collector.
func New(producers []Producer, consumers []Consumer) *Collector {
	return &Collector{
		s:         newStore(),
		producers: producers,
		consumers: consumers,
	}
}

func (c *Collector) Collect() error {
	diffs, err := c.produce()
	if err != nil {
		return errors.Wrap(err, "produce")
	}

	if err := c.consume(diffs); err != nil {
		return errors.Wrap(err, "consume")
	}

	return nil
}

func (c *Collector) produce() ([]*node.Node, error) {
	setCalls := make([]*node.Node, 0)
	c.s.setCallback = func(n *node.Node) {
		setCalls = append(setCalls, n)
	}

	for i, p := range c.producers {
		if err := p.Produce(c.s); err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("producer #%d", i))
		}
	}

	return setCalls, nil
}

func (collector *Collector) consume(setCalls []*node.Node) error {
	var from, to int
	consumerSetCalls := make(map[int]Consumer)
	for {
		from = to
		to = len(setCalls)
		if from == to {
			break
		} else {
			log.Debugf("consuming setCalls %d:%d", from, to)
		}

		for ; from < to; from++ {
			log.Debugf("consuming setCall %s", setCalls[from])
			for i, c := range collector.consumers {
				if consumerSetCalls[from] == c {
					continue
				}

				collector.s.setCallback = func(n *node.Node) {
					setCalls = append(setCalls, n)
					consumerSetCalls[len(setCalls)-1] = c
				}

				setCall := setCalls[from]
				if err := c.Consume(collector.s, setCall); err != nil {
					return errors.Wrap(err, fmt.Sprintf("consumer #%d, setCall %s", i, setCall))
				}
			}
		}
	}

	return nil
}
