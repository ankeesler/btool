package unmarshaler_test

import (
	"testing"

	nodev1 "github.com/ankeesler/btool/node/api/v1"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalerUnmarshal(t *testing.T) {
	//	//data := []struct {
	//	//	sFunc func() collector.Store
	//	//	from  *nodev1.Node
	//	//	root  string
	//	//
	//	//	to  *node.Node
	//	//	err error
	//	//}{
	//	//	{
	//	//		sFunc: func() collector.Store {
	//	//			b := node.New("b")
	//	//			c := node.New("c")
	//	//			return testutil.FakeStore(b, c)
	//	//		},
	//	//		from: &nodev1.Node{
	//	//			Name:         "a",
	//	//			Dependencies: []string{"b", "c"},
	//	//		},
	//	//		root: "some/root",
	//	//
	//	//		to: &node.Node{
	//	//			Name: "some/root/a",
	//	//			Dependencies:
	//	//		},
	//	//	},
	//	//}
	//
	//	u := unmarshaler.New()
	//
	//	b := node.New("some/root/b")
	//	c := node.New("some/root/c")
	//	s := testutil.FakeStore(b, c)
	//
	//	from := &nodev1.Node{
	//		Name:         "a",
	//		Dependencies: []string{"b", "c"},
	//	}
	//	to := node.New(
	//		"some/root/a",
	//	).Dependency(
	//		b,
	//		c,
	//	).Label(
	//		"io.btool.collector.cc.includePaths",
	//		[]string{
	//			"/some/root/include/dir0",
	//			"/some/root/include/dir1",
	//		},
	//	)
	//
	//	n, err := u.Unmarshal(s, from, "some/root")
	//	require.Nil(t, err)
	//	assert.Equal(t, to, n)
}

func marshalAny(t *testing.T, strings ...string) *any.Any {
	m := nodev1.Strings{
		Strings: strings,
	}
	any, err := ptypes.MarshalAny(&m)
	require.Nil(t, err)
	return any
}
