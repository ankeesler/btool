// Package store provides an abstraction of registry and project storage
// directories.
package store

import "path/filepath"

// Store provides registry and project directories as children of a cache
// directory.
type Store struct {
	cache string
}

// New creates a new Store with a cache directory as its root.
func New(cache string) *Store {
	return &Store{
		cache: cache,
	}
}

func (s *Store) RegistryDir(registry string) (string, error) {
	return filepath.Join(s.cache, "registries", registry), nil
}

func (s *Store) ProjectDir(project string) (string, error) {
	return filepath.Join(s.cache, "projects", project), nil
}
