// Package prodcon provides a producer/consumer framework that can act as a
// collector.
//
// Producer's add node.Node's to a Store. Consumer's are notified of node.Node's
// being added to the Store. A Store is just a place where node.Node's are kept.
package prodcon

import (
	"fmt"

	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Producer

// Producer is a type that adds node.Node's to a Store.
type Producer interface {
	Produce(*Store) error
}

// Type of the CRUD action on a node.Node.
type DiffType int

const (
	// A node.Node has been added.
	DiffAdd DiffType = iota
)

// Diff describes a change to a node.Node.
type Diff struct {
	Type DiffType

	// Name of the node.Node that has been CRUD'd.
	Name string
}

func (d *Diff) String() string {
	var teyep string
	switch d.Type {
	case DiffAdd:
		teyep = "DiffAdd"
	default:
		teyep = "???"
	}
	return fmt.Sprintf("%s:%s", teyep, d.Name)
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Consumer

// Consumer is a type that reacts to Producer's adding node.Node's to a Store.
type Consumer interface {
	Consume(*Store, *Diff) error
}

// PC (i.e., "prodcon") is an object that will run Producer's and Consumer's.
//
// This particular "prodcon" is very synchronus. It will:
//   1. run each Producer to completion
//   2. run all Consumer's on every Diff from a Producer
//   3. run all Consumer's on every Diff from a Consumer that isn't their own
//   4. repeat 4 until there are no more Diff's left
type PC struct {
	s         *Store
	producers []Producer
	consumers []Consumer
}

// New creates a new PC.
func New(s *Store, producers []Producer, consumers []Consumer) *PC {
	return &PC{
		s:         s,
		producers: producers,
		consumers: consumers,
	}
}

// Run runs the PC.
func (pc *PC) Run() error {
	diffs, err := pc.produce()
	if err != nil {
		return errors.Wrap(err, "produce")
	}

	if err := pc.consume(diffs); err != nil {
		return errors.Wrap(err, "consume")
	}

	return nil
}

func (pc *PC) produce() ([]*Diff, error) {
	diffs := make([]*Diff, 0)
	pc.s.addCallback = func(n *node.Node) {
		diffs = append(diffs, &Diff{
			Type: DiffAdd,
			Name: n.Name,
		})
	}

	for i, p := range pc.producers {
		if err := p.Produce(pc.s); err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("producer #%d", i))
		}
	}

	return diffs, nil
}

func (pc *PC) consume(diffs []*Diff) error {
	var from, to int
	consumerDiffs := make(map[int]Consumer)
	for {
		from = to
		to = len(diffs)
		if from == to {
			break
		} else {
			log.Debugf("consuming diffs %d:%d", from, to)
		}

		for ; from < to; from++ {
			for i, c := range pc.consumers {
				if consumerDiffs[from] == c {
					continue
				}

				pc.s.addCallback = func(n *node.Node) {
					diff := &Diff{
						Type: DiffAdd,
						Name: n.Name,
					}
					diffs = append(diffs, diff)
					consumerDiffs[len(diffs)-1] = c
				}

				diff := diffs[from]
				if err := c.Consume(pc.s, diff); err != nil {
					return errors.Wrap(err, fmt.Sprintf("consumer #%d, diff %s", i, diff))
				}
			}
		}
	}

	return nil
}
