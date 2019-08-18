// Package main provides the btool registry API.
package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node/registry"
	"github.com/ankeesler/btool/node/registry/api"
	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

func main() {
	if err := run(); err != nil {
		log.Errorf(err.Error())
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

	level, err := log.ParseLevel(*loglevel)
	if err != nil {
		return errors.Wrap(err, "parse log level")
	}
	log.CurrentLevel = level
	log.Debugf("log level set to %s", level)

	name := getRegistryName(*address)
	r, err := registry.CreateFSRegistry(afero.NewOsFs(), *dir, name)
	if err != nil {
		return errors.Wrap(err, "create registry")
	}

	log.Infof("listening on %s", *address)
	api := api.New(r)
	if err := http.ListenAndServe(*address, api); err != nil {
		return errors.Wrap(err, "listen and serve")
	}

	return nil
}

func getRegistryName(deefault string) string {
	app, err := cfenv.Current()
	if err != nil {
		log.Debugf("note: cannot get cfenv: %s", err.Error())
	} else if len(app.ApplicationURIs) == 0 {
		log.Debugf("note: no cfenv application uris")
	} else {
		return app.ApplicationURIs[0]
	}

	return deefault
}
