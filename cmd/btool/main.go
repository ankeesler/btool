// Package main provides the btool main function.
package main

import (
	"flag"
	"os"

	"github.com/ankeesler/btool"
	"github.com/ankeesler/btool/formatter"
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

	cfg := btool.Cfg{
		Root:   *root,
		Cache:  *cache,
		Target: *target,

		CompilerC:  "clang",
		CompilerCC: "clang++",
		Linker:     "clang",
	}
	if err := btool.Build(&cfg); err != nil {
		return errors.Wrap(err, "build")
	}

	return nil
}
