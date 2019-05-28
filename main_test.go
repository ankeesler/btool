package main_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestBuildingC(t *testing.T) {
	btool := build(t)

	cache, err := ioutil.TempDir("", "btool_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(cache)

	for i := 0; i < 3; i++ {
		output, err := exec.Command(
			btool,
			"--root",
			"fixture/ComplexC",
			"--cache",
			cache,
			"build",
			"fixture/ComplexC/main.c",
		).CombinedOutput()
		if err != nil {
			t.Fatal(string(output), err)
		}

		if output, err := exec.Command(
			"rm",
			filepath.Join(
				cache,
				"objects",
				fmt.Sprintf("dep-%d/dep-%da.o", i, i),
			),
		).CombinedOutput(); err != nil {
			t.Fatal(string(output), err)
		}
	}

	main := filepath.Join(cache, "binaries", "out")
	if _, err := os.Stat(main); err != nil {
		t.Fatalf("%s does not exist: %s", main, err.Error())
	}

	output, err := exec.Command(main).CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	if string(output) != "hey! i am running!\n" {
		t.Fatal("wrong output")
	}
}

func TestBuildingCC(t *testing.T) {
	btool := build(t)

	cache, err := ioutil.TempDir("", "btool_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(cache)

	for i := 0; i < 3; i++ {
		output, err := exec.Command(
			btool,
			"--root",
			"fixture/ComplexCC",
			"--cache",
			cache,
			"build",
			"fixture/ComplexCC/main.cc",
		).CombinedOutput()
		if err != nil {
			t.Fatal(string(output), err)
		}

		if output, err := exec.Command(
			"rm",
			filepath.Join(
				cache,
				"objects",
				fmt.Sprintf("dep-%d/dep-%da.o", i, i),
			),
		).CombinedOutput(); err != nil {
			t.Fatal(string(output), err)
		}
	}

	main := filepath.Join(cache, "binaries", "out")
	if _, err := os.Stat(main); err != nil {
		t.Fatalf("%s does not exist: %s", main, err.Error())
	}

	output, err := exec.Command(main).CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	if string(output) != "hey! i am running!\n" {
		t.Fatal("wrong output")
	}
}

func TestScanningC(t *testing.T) {
	btool := build(t)

	cache, err := ioutil.TempDir("", "btool_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(cache)

	output0, err := exec.Command(
		btool,
		"--root",
		"fixture/BasicC",
		"--cache",
		cache,
		"scan",
	).CombinedOutput()
	if err != nil {
		t.Fatal(string(output0), err)
	}

	output1, err := exec.Command(
		btool,
		"--root",
		"fixture/BasicC",
		"--cache",
		cache,
		"scan",
		"fixture/BasicC/main.c",
	).CombinedOutput()
	if err != nil {
		t.Fatal(string(output1), err)
	}
}

func TestScanningCC(t *testing.T) {
	btool := build(t)

	cache, err := ioutil.TempDir("", "btool_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(cache)

	output0, err := exec.Command(
		btool,
		"--root",
		"fixture/BasicCC",
		"--cache",
		cache,
		"scan",
	).CombinedOutput()
	if err != nil {
		t.Fatal(string(output0), err)
	}

	output1, err := exec.Command(
		btool,
		"--root",
		"fixture/BasicCC",
		"--cache",
		cache,
		"scan",
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
