package main_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestBuildingC(t *testing.T) {
	btool := build(t)

	store, err := ioutil.TempDir("", "btool_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(store)

	output, err := exec.Command(
		btool,
		"-root",
		"fixture/ComplexC",
		"-store",
		store,
		"-target",
		"fixture/ComplexC/main.c",
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

func TestBuildingCC(t *testing.T) {
	btool := build(t)

	store, err := ioutil.TempDir("", "btool_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(store)

	output, err := exec.Command(
		btool,
		"-root",
		"fixture/ComplexCC",
		"-store",
		store,
		"-target",
		"fixture/ComplexCC/main.cc",
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

func TestScanningC(t *testing.T) {
	btool := build(t)

	store, err := ioutil.TempDir("", "btool_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(store)

	output0, err := exec.Command(
		btool,
		"-root",
		"fixture/BasicC",
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
		"fixture/BasicC",
		"-store",
		store,
		"-scan",
		"-target",
		"fixture/BasicC/main.c",
	).CombinedOutput()
	if err != nil {
		t.Fatal(string(output1), err)
	}
}

func TestScanningCC(t *testing.T) {
	btool := build(t)

	store, err := ioutil.TempDir("", "btool_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(store)

	output0, err := exec.Command(
		btool,
		"-root",
		"fixture/BasicCC",
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
		"fixture/BasicCC",
		"-store",
		store,
		"-scan",
		"-target",
		"fixture/BasicCC/main.cc",
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
