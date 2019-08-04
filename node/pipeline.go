package node

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Handler

type Config struct {
	Root   string
	Cache  string
	Target string

	CCompiler  string
	CCCompiler string
}

type Handler interface {
	Handle(*Config, []*Node) ([]*Node, error)
}

func Pipeline(cfg *Config, handlers ...Handler) error {
	nodes := make([]*Node, 0)
	for _, handler := range handlers {
		var err error
		nodes, err = handler.Handle(cfg, nodes)
		if err != nil {
			return err
		}
	}
	return nil
}
