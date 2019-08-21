// Package main provides the btool main function.
package main

import (
	"flag"
	"os"
	"runtime"
	"strings"

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
	loglevel := flag.String("loglevel", "info", "Verbosity of log")
	root := flag.String("root", ".", "Root of node list")
	cache := flag.String("cache", ".btool", "Cache directory")
	target := flag.String("target", "main", "Target to build")
	registries := flag.String(
		"registries",
		"https://btoolregistry.cfapps.io",
		"List of registries, e.g., https://a.io,file://path/to/reg/dir",
	)
	dryrun := flag.Bool("dryrun", false, "Do not actually build")
	version := flag.Bool("version", false, "Print version")
	help := flag.Bool("help", false, "Show this help message")

	flag.Parse()
	if *help {
		flag.Usage()
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

		CompilerC:  tc.CompilerC,
		CompilerCC: tc.CompilerCC,
		Archiver:   tc.Archiver,
		LinkerC:    tc.LinkerC,
		LinkerCC:   tc.LinkerCC,

		Registries: strings.Split(*registries, ","),
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
