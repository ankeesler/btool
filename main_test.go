package main_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestBuilding(t *testing.T) {
	btool := build(t)

	store, err := ioutil.TempDir("", "btool_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(store)

	output, err := exec.Command(
		btool,
		"-root",
		"fixture/Complex",
		"-store",
		store,
		"-target",
		"fixture/Complex/main.c",
	).CombinedOutput()
	if err != nil {
		t.Fatal(string(output), err)
	}

	main := filepath.Join(store, "binaries", "out")
	if _, err := os.Stat(main); err != nil {
		t.Fatalf("%s does not exist: %s", main, err.Error())
	}

	output, err = exec.Command(main).CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	if string(output) != "hey! i am running!\n" {
		t.Fatal("wrong output")
	}
}

func TestScanning(t *testing.T) {
	btool := build(t)

	store, err := ioutil.TempDir("", "btool_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(store)

	output0, err := exec.Command(
		btool,
		"-root",
		"fixture/Basic",
		"-store",
		store,
		"-scan",
	).CombinedOutput()
	if err != nil {
		t.Fatal(string(output0), err)
	}

	output1, err := exec.Command(
		btool,
		"-root",
		"fixture/Basic",
		"-store",
		store,
		"-scan",
		"-target",
		"fixture/Basic/main.c",
	).CombinedOutput()
	if err != nil {
		t.Fatal(string(output1), err)
	}
}

func build(t *testing.T) string {
	btool := "/tmp/btool"
	output, err := exec.Command(
		"go",
		"build",
		"-o",
		btool,
		"github.com/ankeesler/btool",
	).CombinedOutput()
	if err != nil {
		t.Fatal(string(output), err)
	}
	return btool
}
