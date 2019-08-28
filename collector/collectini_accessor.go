package collector

// CollectiniAccessor is a dumb CollectiniCreator. It simple returns the
// Collectini that it was initialized with.
type CollectiniAccessor struct {
	c Collectini
}

// NewCollectiniAccessor create a new CollectiniAccessor.
func NewCollectiniAccessor(c Collectini) *CollectiniAccessor {
	return &CollectiniAccessor{
		c: c,
	}
}

// Create simply returns the Collectini that this CollectiniAccessor was
// initialized with and a nil error.
func (ca *CollectiniAccessor) Create() (Collectini, error) {
	return ca.c, nil
}
