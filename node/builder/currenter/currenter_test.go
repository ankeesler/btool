package currenter_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/builder/currenter"
	"github.com/stretchr/testify/require"
)

func TestCurrenterCurrent(t *testing.T) {
	root, err := ioutil.TempDir("", "btool_current_test")
	require.Nil(t, err)
	defer func() {
		require.Nil(t, os.RemoveAll(root))
	}()

	// a -> b, c
	// b -> c
	// c
	nodeC := node.New(filepath.Join(root, "c"))
	nodeB := node.New(filepath.Join(root, "b")).Dependency(nodeC)
	nodeA := node.New(filepath.Join(root, "a")).Dependency(nodeB, nodeC)
	_ = nodeA

	data := []struct {
		name     string
		n        *node.Node
		written  []*node.Node
		expected bool
	}{
		{
			name:     "NotWritten",
			written:  []*node.Node{},
			n:        nodeC,
			expected: false,
		},
		{
			name:     "Written",
			written:  []*node.Node{nodeC},
			n:        nodeC,
			expected: true,
		},
		{
			name:     "ParentOlder",
			written:  []*node.Node{nodeC, nodeB},
			n:        nodeB,
			expected: true,
		},
		{
			name:     "ParentNewer",
			written:  []*node.Node{nodeB, nodeC},
			n:        nodeB,
			expected: false,
		},
		{
			name:     "AncestorNewer",
			written:  []*node.Node{nodeA, nodeC},
			n:        nodeA,
			expected: false,
		},
	}
	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			for _, w := range datum.written {
				time.Sleep(time.Millisecond * 100)
				require.Nil(
					t,
					ioutil.WriteFile(w.Name, []byte(w.Name), 0644),
				)
			}

			c := currenter.New()
			current, err := c.Current(datum.n)
			require.Nil(t, err)
			require.Equal(t, datum.expected, current)
		})
	}
}
