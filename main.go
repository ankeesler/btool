package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/ankeesler/btool/builder"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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

	store := flag.String("store", filepath.Join(cwd, ".btool"), "Path to btool store")
	target := flag.String("target", "main.c", "Path to build target")
	logLevel := flag.String("loglevel", "info", "Log level")
	help := flag.Bool("help", false, "Print this help message")

	flag.Parse()

	if *help {
		flag.Usage()
		return nil
	}

	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "15:04:05.000000000",
	})

	level, err := log.ParseLevel(*logLevel)
	if err != nil {
		return errors.Wrap(err, "parse log level")
	}
	log.SetLevel(level)

	b := builder.New(*store)
	if err := b.Build(*target); err != nil {
		return errors.Wrap(err, "build")
	}

	return nil
}
