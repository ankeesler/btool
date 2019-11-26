package resolvers_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/resolvers"
	"github.com/stretchr/testify/require"
)

func TestUntar(t *testing.T) {
	outputDir, err := ioutil.TempDir("", "btool_untar_test")
	require.Nil(t, err)
	defer func() {
		require.Nil(t, os.RemoveAll(outputDir))
	}()
	resolversTarGz := filepath.Join(outputDir, "resolvers.tar.gz")

	cmd := exec.Command("tar")
	cmd.Args = append(cmd.Args, "czf")
	cmd.Args = append(cmd.Args, resolversTarGz)
	cmd.Args = append(cmd.Args, "untar_test.go")
	output, err := cmd.CombinedOutput()
	require.Nil(t, err, string(output))

	u := resolvers.NewUntar()
	n := node.New(filepath.Join(outputDir, "untar_test.go"))
	d := node.New(resolversTarGz)
	n.Dependency(d)
	require.Nil(t, u.Resolve(n))

	_, err = os.Stat(n.Name)
	require.Nil(t, err)
}
