package integration_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestBuild(t *testing.T) {
	genfixture, err := build("github.com/ankeesler/btool/node/cmd/genfixture")
	if err != nil {
		t.Fatal(err)
	}

	build, err := build("github.com/ankeesler/btool/node/cmd/build")
	if err != nil {
		t.Fatal(err)
	}

	root := filepath.Join(os.TempDir(), "btool_node_integration_test")
	os.RemoveAll(root)
	os.MkdirAll(root, 0755)
	if err := exec.Command(genfixture, "-root", root).Run(); err != nil {
		t.Fatal(err)
	}

	t.Run("Object", func(t *testing.T) {
		if err := exec.Command(
			build,
			"-target",
			"main.o",
			"-root",
			root,
		).Run(); err != nil {
			t.Error(err)
		}
	})
}

func testObject(t *testing.T) {

}

func build(path string) (string, error) {
	name := filepath.Join(os.TempDir(), filepath.Base(path))
	if err := exec.Command(
		"go",
		"build",
		"-o",
		name,
		path,
	).Run(); err != nil {
		return "", err
	}

	return name, nil
}
