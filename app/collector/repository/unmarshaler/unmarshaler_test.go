package unmarshaler_test

import (
	"testing"

	"github.com/ankeesler/btool/app/collector/repository/unmarshaler"
	"github.com/ankeesler/btool/app/collector/testutil"
	"github.com/ankeesler/btool/node"
	nodev1 "github.com/ankeesler/btool/node/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalerUnmarshal(t *testing.T) {
	//data := []struct {
	//	sFunc func() collector.Store
	//	from  *nodev1.Node
	//	root  string
	//
	//	to  *node.Node
	//	err error
	//}{
	//	{
	//		sFunc: func() collector.Store {
	//			b := node.New("b")
	//			c := node.New("c")
	//			return testutil.FakeStore(b, c)
	//		},
	//		from: &nodev1.Node{
	//			Name:         "a",
	//			Dependencies: []string{"b", "c"},
	//		},
	//		root: "some/root",
	//
	//		to: &node.Node{
	//			Name: "some/root/a",
	//			Dependencies:
	//		},
	//	},
	//}

	u := unmarshaler.New()

	b := node.New("some/root/b")
	c := node.New("some/root/c")
	s := testutil.FakeStore(b, c)

	from := &nodev1.Node{
		Name:         "a",
		Dependencies: []string{"b", "c"},
	}
	to := node.New("some/root/a").Dependency(b, c)

	n, err := u.Unmarshal(s, from, "some/root")
	require.Nil(t, err)
	assert.Equal(t, to, n)
}
