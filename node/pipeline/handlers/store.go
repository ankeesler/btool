package handlers

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Store

// Store is an abstraction of places that these handlers store node.Node's.
type Store interface {
	ProjectDir(project string) (string, error)
	RegistryDir(registry string) (string, error)
}
