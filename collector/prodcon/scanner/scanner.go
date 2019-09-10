// Package scanner provides a prodcon.Producer that can produce node.Node's by
// walking a filesystem.
package scanner

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Collectini

// Walker is a type that can walk a filesystem and return a list of paths.
type Walker interface {
	Walk(string, []string) ([]string, error)
}

// Scanner is a type that can Produce() node.Node's from a filesystem.
type Scanner struct {
	w Walker
}

func New(w Walker) *Scanner {
	return &Scanner{
		w: w,
	}
}

func (s *Scanner) Produce() error {
	return nil
}
