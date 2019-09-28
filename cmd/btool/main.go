// Package main provides the btool main function.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ankeesler/btool"
	"github.com/ankeesler/btool/log"
	"github.com/pkg/errors"
)

const v = "0.0.1"

func main() {
	if err := run(); err != nil {
		log.Errorf(err.Error())
		os.Exit(1)
	}
}

func run() error {
	home, ok := os.LookupEnv("HOME")
	if !ok {
		home = "."
	}

	loglevel := flag.String("loglevel", "info", "Set log level (debug, info, error)")
	root := flag.String("root", ".", "Root of node list")
	cache := flag.String("cache", filepath.Join(home, ".btool"), "Cache directory")
	target := flag.String("target", "main", "Target to build")
	registry := flag.String(
		"registry",
		"https://btoolregistry.cfapps.io",
		"Btool registry link, e.g., file://path/to/reg/dir or https://a.io",
	)
	dryrun := flag.Bool("dryrun", false, "Do not actually build (or run)")
	clean := flag.Bool("clean", false, "Clean all nodes")
	list := flag.Bool("list", false, "Simply list all targets")
	run := flag.Bool("run", false, "Run the target (after building)")
	watch := flag.Bool("watch", false, "Watch for file changes and rebuild (or rerun)")
	version := flag.Bool("version", false, "Print version")
	help := flag.Bool("help", false, "Show this help message")

	flag.Parse()
	if *help {
		printHelp()
		os.Exit(1)
	}

	if *version {
		log.Infof("version %s", v)
		return nil
	}

	level, err := log.ParseLevel(*loglevel)
	if err != nil {
		return errors.Wrap(err, "parse log level")
	}
	log.CurrentLevel = level

	tc, err := findToolchain()
	if err != nil {
		return errors.Wrap(err, "find toolchain")
	}

	cfg := btool.Cfg{
		Root:   *root,
		Cache:  *cache,
		Target: *target,

		DryRun: *dryrun,
		Clean:  *clean,
		List:   *list,
		Run:    *run,
		Watch:  *watch,

		CompilerC: tc.CompilerC,
		CompilerCFlags: []string{
			"-Wall",
			"-Werror",
			"-g",
			"-O0",
			"--std=c17",
		},
		CompilerCC: tc.CompilerCC,
		CompilerCCFlags: []string{
			"-Wall",
			"-Werror",
			"-g",
			"-O0",
			"--std=c++17",
		},
		Archiver: tc.Archiver,
		LinkerC:  tc.LinkerC,
		LinkerCC: tc.LinkerCC,

		Registry: *registry,
	}
	if err := btool.Run(&cfg); err != nil {
		return errors.Wrap(err, "run")
	}

	return nil
}

type toolchain struct {
	CompilerC  string
	CompilerCC string
	Archiver   string
	LinkerC    string
	LinkerCC   string
}

var (
	linux = toolchain{
		CompilerC:  "gcc",
		CompilerCC: "g++",
		Archiver:   "ar",
		LinkerC:    "gcc",
		LinkerCC:   "g++",
	}
	darwin = toolchain{
		CompilerC:  "clang",
		CompilerCC: "clang++",
		Archiver:   "ar",
		LinkerC:    "clang",
		LinkerCC:   "clang++",
	}
)

func findToolchain() (*toolchain, error) {
	switch runtime.GOOS {
	case "linux":
		return &linux, nil
	case "darwin":
		return &darwin, nil
	default:
		return nil, errors.New("unknown toolchain for system")
	}
}

func printHelp() {
	fmt.Println(os.Args[0], "- The simplest C/C++ build tool.")

	fmt.Println()
	fmt.Println("Examples")
	fmt.Println("  btool -target main -root some/root")
	fmt.Println("  \tbuild an executable main from project root at ./some/root")
	fmt.Println("  btool -target main -run")
	fmt.Println("  \tbuild and run an executable main")
	fmt.Println("  btool -target main -run -watch")
	fmt.Println("  \tbuild and run an executable main continuously on file changes")
	fmt.Println("  btool -list -root some/root")
	fmt.Println("  \tlist the possible targets in project ./some/root")
	fmt.Println("  btool -clean -target main")
	fmt.Println("  \tdelete all built files associated with target main")

	fmt.Println()
	fmt.Println("Flags")
	flag.PrintDefaults()
}
