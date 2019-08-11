// Package main provides the btool main function.
package main

import (
	"flag"
	"os"
	"strings"

	"github.com/ankeesler/btool"
	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/toolchain"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(formatter.New())
	if err := run(); err != nil {
		logrus.Error(err)
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
		"List of registries (e.g., https://a.io,file://path/to/reg/dir)",
	)
	help := flag.Bool("help", false, "Show this help message")

	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(1)
	}

	level, err := logrus.ParseLevel(*loglevel)
	if err != nil {
		return errors.Wrap(err, "parse log level")
	}
	logrus.SetLevel(level)

	tc, err := toolchain.Find()
	if err != nil {
		return errors.Wrap(err, "toolchain find")
	}

	cfg := btool.Cfg{
		Root:   *root,
		Cache:  *cache,
		Target: *target,

		CompilerC:  tc.CompilerC,
		CompilerCC: tc.CompilerCC,
		Linker:     tc.Linker,

		Registries: strings.Split(*registries, ","),
	}
	if err := btool.Run(&cfg); err != nil {
		return errors.Wrap(err, "run")
	}

	return nil
}
