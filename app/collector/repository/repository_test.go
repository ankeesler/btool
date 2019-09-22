package repository_test

import (
	"path/filepath"
	"testing"

	"github.com/ankeesler/btool/app/collector/repository"
	"github.com/ankeesler/btool/app/collector/repository/repositoryfakes"
	"github.com/ankeesler/btool/app/collector/testutil"
	"github.com/ankeesler/btool/node"
	nodev1 "github.com/ankeesler/btool/node/api/v1"
	"github.com/ankeesler/btool/node/api/v1/v1fakes"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepositoryProduce(t *testing.T) {
	cache := "some/cache"

	fs := afero.NewMemMapFs()

	repoANodes := []*nodev1.Node{
		// a -> b
		// b -> c
		// c
		&nodev1.Node{
			Name:         "c",
			Dependencies: []string{""},
		},
		&nodev1.Node{
			Name:         "b",
			Dependencies: []string{"c"},
		},
		&nodev1.Node{
			Name:         "a",
			Dependencies: []string{"b"},
		},
	}
	repoARev := "abc123"

	repoBNodes := []*nodev1.Node{
		// d -> e
		// e -> f
		// f
		&nodev1.Node{
			Name:         "f",
			Dependencies: []string{""},
		},
		&nodev1.Node{
			Name:         "e",
			Dependencies: []string{"f"},
		},
		&nodev1.Node{
			Name:         "d",
			Dependencies: []string{"e"},
		},
	}
	repoBRev := "def456"

	c := &v1fakes.FakeRegistryClient{}
	c.ListRepositoriesReturnsOnCall(0, &nodev1.ListRepositoriesResponse{
		Repositories: []*nodev1.Repository{
			&nodev1.Repository{
				Name:     "repo-a",
				Revision: repoARev,
				Nodes:    repoANodes,
			},

			&nodev1.Repository{
				Name:     "repo-b",
				Revision: repoBRev,
				Nodes:    repoBNodes,
			},
		},
	}, nil)

	nodeC := &node.Node{
		Name:         "abc123/c",
		Dependencies: []*node.Node{},
	}
	nodeB := &node.Node{
		Name:         "abc123/b",
		Dependencies: []*node.Node{nodeC},
	}
	nodeA := &node.Node{
		Name:         "abc123/a",
		Dependencies: []*node.Node{nodeB},
	}

	nodeF := &node.Node{
		Name:         "def456/f",
		Dependencies: []*node.Node{},
	}
	nodeE := &node.Node{
		Name:         "def456/e",
		Dependencies: []*node.Node{nodeF},
	}
	nodeD := &node.Node{
		Name:         "def456/d",
		Dependencies: []*node.Node{nodeE},
	}

	setNodes := []*node.Node{
		nodeC,
		nodeB,
		nodeA,

		nodeF,
		nodeE,
		nodeD,
	}
	u := &repositoryfakes.FakeUnmarshaler{}
	for i, n := range setNodes {
		u.UnmarshalReturnsOnCall(i, n, nil)
	}

	r := repository.New(fs, c, u, cache)

	s := testutil.FakeStore()
	require.Nil(t, r.Produce(s))

	assert.Equal(t, 1, c.ListRepositoriesCallCount())

	assert.Equal(t, 6, u.UnmarshalCallCount())
	for i := 0; i < u.UnmarshalCallCount(); i++ {
		acS, acN, acRev := u.UnmarshalArgsForCall(i)
		assert.Equal(t, s, acS)
		if i < 3 {
			assert.Equal(t, repoANodes[i], acN)
			assert.Equal(t, acRev, filepath.Join(cache, repoARev))
		} else {
			assert.Equal(t, repoBNodes[i-3], acN)
			assert.Equal(t, acRev, filepath.Join(cache, repoBRev))
		}
	}

	assert.Equal(t, 6, s.SetCallCount())
	for i := 0; i < s.SetCallCount(); i++ {
		assert.Equal(t, setNodes[i], s.SetArgsForCall(i))
	}
}
