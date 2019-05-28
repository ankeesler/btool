// Package config provides an abstraction for per-project settings.
package config

type Config struct {
	Name  string
	Root  string
	Cache string
}

func New(name, root, cache string) *Config {
	return &Config{
		Name:  name,
		Root:  root,
		Cache: cache,
	}
}
