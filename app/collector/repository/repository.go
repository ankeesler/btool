// Package repository provides a collector.Producer that gets node.Node's from a
// NodeRepository API.
package repository

import (
	"context"
	"path/filepath"
	"time"

	"github.com/ankeesler/btool/app/collector"
	"github.com/ankeesler/btool/node"
	nodev1 "github.com/ankeesler/btool/node/api/v1"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Unmarshaler

// Unmarshaler is a type that, given a Node and a root for the Node, can produce
// a node.Node.
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
	repositories, err := r.listRepositories()
	if err != nil {
		return errors.Wrap(err, "list repositories")
	}

	for _, repository := range repositories {
		root := filepath.Join(r.cache, repository.Revision)

		// TODO: caching?

		for _, repositoryNode := range repository.Nodes {
			n, err := r.u.Unmarshal(s, repositoryNode, root)
			if err != nil {
				return errors.Wrap(err, "unmarshal")
			}

			s.Set(n)
		}
	}

	return nil
}

func (r *Repository) listRepositories() ([]*nodev1.Repository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	req := nodev1.ListRepositoriesRequest{}

	rsp, err := r.c.ListRepositories(ctx, &req)
	if err != nil {
		return nil, errors.Wrap(err, "list repositories")
	}

	return rsp.Repositories, nil
}
