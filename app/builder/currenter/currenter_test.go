package currenter_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ankeesler/btool/app/builder/currenter"
	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
	"github.com/stretchr/testify/require"
)

func TestCurrenterCurrent(t *testing.T) {
	data := []struct {
		name     string
		n        string
		written  []string
		expected bool
	}{
		{
			name:     "NotWritten",
			written:  []string{},
			n:        "c",
			expected: false,
		},
		{
			name:     "Written",
			written:  []string{"c"},
			n:        "c",
			expected: true,
		},
		{
			name:     "ParentOlder",
			written:  []string{"c", "b"},
			n:        "b",
			expected: true,
		},
		{
			name:     "ParentNewer",
			written:  []string{"b", "c"},
			n:        "b",
			expected: false,
		},
		{
			name:     "AncestorNewer",
			written:  []string{"e", "d", "f"},
			n:        "d",
			expected: false,
		},
	}
	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			root, err := ioutil.TempDir("", "btool_current_test")
			require.Nil(t, err)
			defer func() {
				require.Nil(t, os.RemoveAll(root))
			}()

			// a -> b, c
			// b -> c
			// c
			// d -> e
			// e -> f
			// f
			nodeF := node.New(filepath.Join(root, "f"))
			nodeE := node.New(filepath.Join(root, "e")).Dependency(nodeF)
			nodeD := node.New(filepath.Join(root, "d")).Dependency(nodeE)
			nodeC := node.New(filepath.Join(root, "c"))
			nodeB := node.New(filepath.Join(root, "b")).Dependency(nodeC)
			nodeA := node.New(filepath.Join(root, "a")).Dependency(nodeB, nodeC)

			n := find(datum.n, nodeA, nodeB, nodeC, nodeD, nodeE, nodeF)
			require.NotNil(t, n)

			for _, w := range datum.written {
				file := filepath.Join(root, w)
				time.Sleep(time.Millisecond * 100)
				require.Nil(
					t,
					ioutil.WriteFile(file, []byte(file), 0644),
				)
				info, err := os.Stat(file)
				require.Nil(t, err)
				log.Debugf("%s mod time %s", file, info.ModTime())
			}

			c := currenter.New()
			current, err := c.Current(n)
			require.Nil(t, err)
			require.Equal(t, datum.expected, current)
		})
	}
}

func find(name string, nodes ...*node.Node) *node.Node {
	for _, n := range nodes {
		if strings.HasSuffix(n.Name, name) {
			return n
		}
	}
	return nil
}
