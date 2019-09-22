package repository_test

import (
	"testing"

	"github.com/ankeesler/btool/app/collector/repository"
	"github.com/ankeesler/btool/app/collector/repository/repositoryfakes"
	"github.com/ankeesler/btool/app/collector/testutil"
	"github.com/ankeesler/btool/node/api/v1/v1fakes"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestRepositoryProduce(t *testing.T) {
	fs := afero.NewMemMapFs()
	c := &v1fakes.FakeNodeRepositoryClient{}
	u := &repositoryfakes.FakeUnmarshaler{}
	r := repository.New(fs, c, u, "some/cache")

	s := testutil.FakeStore()
	require.Nil(t, r.Produce(s))
}
