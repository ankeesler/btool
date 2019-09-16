package collector

import "github.com/ankeesler/btool/node"

// TrivialProducer is a type that very simply Produce()'s a provided node.Node.
type TrivialProducer struct {
	n *node.Node
}

// NewTrivialProducer creates a TrivialProducer with a provided node.Node.
func NewTrivialProducer(n *node.Node) *TrivialProducer {
	return &TrivialProducer{
		n: n,
	}
}

func (tp *TrivialProducer) Produce(s Store) error {
	s.Set(tp.n)
	return nil
}
