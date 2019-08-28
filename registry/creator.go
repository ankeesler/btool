package registry

type Creator struct {
	url string
}

func NewCreator(url string) *Creator {
	return &Creator{
		url: url,
	}
}

func (c *Creator) Create() (Registry, error) {
	return nil, nil
}
