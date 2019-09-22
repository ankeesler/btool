// Package repository provides a collector.Producer that gets node.Node's from a
// NodeRepository API.
package repository

import (
	"github.com/ankeesler/btool/app/collector"
	"github.com/ankeesler/btool/node"
	nodev1 "github.com/ankeesler/btool/node/api/v1"
	"github.com/spf13/afero"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Unmarshaler

// Unmarshaler is a type that, given a Node and its Repository's revision from
// the NodeRepository API, can produce a node.Node.
type Unmarshaler interface {
	Unmarshal(collector.Store, *nodev1.Node, string) (*node.Node, error)
}

// Repository is a collector.Producer that can add node.Node's from a
// NodeRepository API to a collector.Store. It does this by translating them
// via an Unmarshaler interface.
type Repository struct {
	fs    afero.Fs
	c     nodev1.RegistryClient
	u     Unmarshaler
	cache string
}

// New creates a new Repository.
func New(
	fs afero.Fs,
	c nodev1.RegistryClient,
	u Unmarshaler,
	cache string,
) *Repository {
	return &Repository{
		fs:    fs,
		c:     c,
		u:     u,
		cache: cache,
	}
}

func (r *Repository) Produce(s collector.Store) error {
	return nil
}
