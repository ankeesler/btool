package integration_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestBuild(t *testing.T) {
	genfixture, err := build("github.com/ankeesler/btool/cmd/genfixture")
	if err != nil {
		t.Fatal(err)
	}

	btool, err := build("github.com/ankeesler/btool/cmd/btool")
	if err != nil {
		t.Fatal(err)
	}

	tmpDir, err := ioutil.TempDir("", "btool_node_integration_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	if err := exec.Command(genfixture, "-root", tmpDir).Run(); err != nil {
		t.Fatal(err)
	}

	t.Log("tmpDir:", tmpDir)

	wd := filepath.Join(tmpDir, "wd")
	if err := os.Mkdir(wd, 0755); err != nil {
		t.Fatal(err)
	}

	t.Run("Object", func(t *testing.T) {
		names := []string{"BasicC", "BasicCC"}
		for _, name := range names {
			t.Run(name, func(t *testing.T) {
				root := filepath.Join(tmpDir, name)
				cache := filepath.Join(tmpDir, "cache")

				cmd := exec.Command(
					btool,
					"-target",
					"dep-0/dep-0.o",
					"-root",
					root,
					"-cache",
					cache,
				)
				cmd.Dir = wd
				if output, err := cmd.CombinedOutput(); err != nil {
					t.Fatal(err, ":", string(output))
				}
			})
		}
	})

	t.Run("Executable", func(t *testing.T) {
		names := []string{"BasicC", "BasicCC"}
		for _, name := range names {
			t.Run(name, func(t *testing.T) {
				root := filepath.Join(tmpDir, name)
				cache := filepath.Join(tmpDir, "cache")

				cmd := exec.Command(
					btool,
					"-target",
					"main",
					"-root",
					root,
					"-cache",
					cache,
				)
				cmd.Dir = wd
				if output, err := cmd.CombinedOutput(); err != nil {
					t.Fatal(err, ":", string(output))
				}

				if err := exec.Command(
					filepath.Join(cache, filepath.Base(root), "main"),
				).Run(); err != nil {
					t.Error(err)
				}
			})
		}
	})
}

func build(path string) (string, error) {
	name := filepath.Join(os.TempDir(), filepath.Base(path))
	cmd := exec.Command(
		"go",
		"build",
		"-o",
		name,
		path,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return name, nil
}
