package node

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Handler

type Handler interface {
	Handle([]*Node) ([]*Node, error)
}

func Pipeline(handlers ...Handler) error {
	nodes := make([]*Node, 0)
	for _, handler := range handlers {
		var err error
		nodes, err = handler.Handle(nodes)
		if err != nil {
			return err
		}
	}
	return nil
}
