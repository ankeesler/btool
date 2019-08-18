// Package integration provides an integration test framework for btool.
package integration

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO: why do main and main.o pop up in the integration directory?
// TODO: why does running btool to create object dep-0/dep-0.o fail?
// TODO: write test for symlink issue (it is in a stash).

// Run runs the integration tests.
func Run(t *testing.T) {
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

	exampleDirs := getExampleDirs(t)

	testCases := []testCase{
		{
			name:     "Object",
			testFunc: object,
		},
		{
			name:     "Executable",
			testFunc: executable,
		},
		{
			name:     "ExecutableLocalRegistry",
			testFunc: executableLocalRegistry,
		},
		{
			name:     "ExecutableRunTwice",
			testFunc: executableRunTwice,
		},
		{
			name:     "Googletest",
			testFunc: googletest,
		},
	}

	for _, testCase := range testCases {
		testCaseTmpDir := filepath.Join(tmpDir, testCase.name)
		t.Run(testCase.name, func(t *testing.T) {
			for _, exampleDir := range exampleDirs {
				example := filepath.Base(exampleDir)
				t.Run(example, func(t *testing.T) {
					testCaseExampleTmpDir := filepath.Join(testCaseTmpDir, example)

					wd := filepath.Join(testCaseExampleTmpDir, "wd")
					if err := os.MkdirAll(wd, 0755); err != nil {
						t.Fatal(err)
					}

					testCase.testFunc(&config{
						btool:   btool,
						root:    exampleDir,
						cache:   filepath.Join(testCaseExampleTmpDir, "cache"),
						wd:      wd,
						example: example,

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

func getExampleDirs(t *testing.T) []string {
	wd, err := os.Getwd()
	require.Nil(t, err)

	examplesDir := filepath.Join(wd, "..", "example")
	infos, err := ioutil.ReadDir(examplesDir)
	require.Nil(t, err)

	exampleDirs := make([]string, 0)
	for _, info := range infos {
		if info.IsDir() {
			exampleDir := filepath.Join(examplesDir, info.Name())
			exampleDirs = append(exampleDirs, exampleDir)
		}
	}

	return exampleDirs
}
