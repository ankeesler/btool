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
	"github.com/ankeesler/btool/scanner/graph"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func main() {
	if err := run(); err != nil {
		logrus.Fatal(err)
	}
}

func run() error {
	cwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "getwd")
	}

	root := flag.String("root", cwd, "Path to project root")
	store := flag.String("store", filepath.Join(cwd, ".btool"), "Path to btool store")
	target := flag.String("target", "", "The file to build")
	scan := flag.Bool("scan", false, "Only perform scan and print out graph")
	logLevel := flag.String("loglevel", "info", "Log level")
	help := flag.Bool("help", false, "Print this help message")

	flag.Parse()

	if *help {
		flag.Usage()
		return nil
	}

	logrus.SetFormatter(formatter.New())

	level, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		return errors.Wrap(err, "parse log level")
	}
	logrus.SetLevel(level)

	fs := afero.NewOsFs()

	var g *graph.Graph
	s := scanner.New(fs, *root)
	if *target == "" {
		g, err = s.ScanRoot()
	} else {
		g, err = s.ScanFile(*target)
	}
	if err != nil {
		return errors.Wrap(err, "scan")
	}

	if *scan {
		logrus.Infof("graph:\n%s", g)
		return nil
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
