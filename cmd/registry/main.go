// Package main provides the btool registry API.
package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/registryapi"
	"github.com/ankeesler/btool/registryapi/registry"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
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
	address := flag.String(
		"address",
		"127.0.0.1:12345",
		"Address on which to listen",
	)
	dir := flag.String(
		"dir",
		".",
		"Directory for registry data",
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

	r, err := registry.Create(afero.NewOsFs(), *dir)
	if err != nil {
		return errors.Wrap(err, "create registry")
	}

	logrus.Infof("listening on %s", *address)
	api := registryapi.New(r)
	if err := http.ListenAndServe(*address, api); err != nil {
		return errors.Wrap(err, "listen and serve")
	}

	return nil
}
