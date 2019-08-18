// Package main provides the btool main function.
package main

import (
	"flag"
	"os"
	"strings"

	"github.com/ankeesler/btool"
	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/toolchain"
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
	output := flag.String("output", *target, "Output file")
	registries := flag.String(
		"registries",
		"https://btoolregistry.cfapps.io",
		"List of registries, e.g., https://a.io,file://path/to/reg/dir",
	)
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

	tc, err := toolchain.Find()
	if err != nil {
		return errors.Wrap(err, "toolchain find")
	}

	cfg := btool.Cfg{
		Root:   *root,
		Cache:  *cache,
		Target: *target,
		Output: *output,

		CompilerC:  tc.CompilerC,
		CompilerCC: tc.CompilerCC,
		Archiver:   tc.Archiver,
		Linker:     tc.Linker,

		Registries: strings.Split(*registries, ","),
	}
	if err := btool.Run(&cfg); err != nil {
		return errors.Wrap(err, "run")
	}

	return nil
}
