package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/ankeesler/btool/builder"
	"github.com/ankeesler/btool/builder/compiler"
	"github.com/ankeesler/btool/builder/linker"
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

	root := flag.String("root", cwd, "Path to project root")
	store := flag.String("store", filepath.Join(cwd, ".btool"), "Path to btool store")
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

	s := scanner.New(fs, *root)
	g, err := s.Scan()
	if err != nil {
		return errors.Wrap(err, "scan")
	}

	b := builder.New(
		fs,
		*root,
		*store,
		compiler.New(),
		linker.New(),
	)
	if err := b.Build(g); err != nil {
		return errors.Wrap(err, "build")
	}

	return nil
}
