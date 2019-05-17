package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/ankeesler/btool/builder"
	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/scanner"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "getwd")
	}

	root := flag.Store("root", cwd, "Path to project root")
	store := flag.String("store", filepath.Join(cwd, ".btool"), "Path to btool store")
	target := flag.String("target", "main.c", "Path to build target")
	logLevel := flag.String("loglevel", "info", "Log level")
	help := flag.Bool("help", false, "Print this help message")

	flag.Parse()

	if *help {
		flag.Usage()
		return nil
	}

	log.SetFormatter(formatter.New())

	level, err := log.ParseLevel(*logLevel)
	if err != nil {
		return errors.Wrap(err, "parse log level")
	}
	log.SetLevel(level)

	fs := afero.NewOsFs()

	g, err := scanner.New(fs, root)
	if err != nil {
		return errors.Wrap(err, "scan")
	}

	b := builder.New(*store)
	if err := b.Build(g); err != nil {
		return errors.Wrap(err, "build")
	}

	return nil
}
