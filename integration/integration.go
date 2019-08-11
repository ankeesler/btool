// Package integration provides an integration test framework for btool.
package integration

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TODO: why do main and main.o pop up in the integration directory?
// TODO: why does running btool to create object dep-0/dep-0.o fail?
// TODO: write test for symlink issue (it is in a stash).

// Run runs the integration tests.
func Run(t *testing.T) {
	genfixture, err := build("github.com/ankeesler/btool/cmd/genfixture")
	if err != nil {
		t.Fatal(err)
	}

	btool, err := build("github.com/ankeesler/btool/cmd/btool")
	if err != nil {
		t.Fatal(err)
	}

	registry, err := build("github.com/ankeesler/btool/cmd/registry")
	if err != nil {
		t.Fatal(err)
	}
	_ = registry

	tmpDir, err := ioutil.TempDir("", "btool_node_integration_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	t.Log("tmpDir:", tmpDir)

	testCases := []testCase{
		{
			name:     "Object",
			testFunc: object,
		},
		{
			name:     "Executable",
			testFunc: executable,
		},
	}

	for _, testCase := range testCases {
		testCaseTmpDir := filepath.Join(tmpDir, testCase.name)
		testCaseFixturesTmpDir := filepath.Join(testCaseTmpDir, "fixtures")
		if err := exec.Command(
			genfixture,
			"-root",
			testCaseFixturesTmpDir,
		).Run(); err != nil {
			t.Fatal(err)
		}

		t.Run(testCase.name, func(t *testing.T) {
			fixtures := []string{"BasicC", "BasicCC"}
			for _, fixture := range fixtures {
				t.Run(fixture, func(t *testing.T) {
					testCaseFixtureTmpDir := filepath.Join(testCaseTmpDir, fixture)

					wd := filepath.Join(testCaseFixtureTmpDir, "wd")
					if err := os.MkdirAll(wd, 0755); err != nil {
						t.Fatal(err)
					}

					testCase.testFunc(&config{
						btool:   btool,
						root:    filepath.Join(testCaseFixturesTmpDir, fixture),
						cache:   filepath.Join(testCaseFixtureTmpDir, "cache"),
						wd:      wd,
						fixture: fixture,

						t: t,
					})
				})
			}
		})
	}
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
