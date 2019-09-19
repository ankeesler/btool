// Package main provides the genclass main function.
package main

import (
	"flag"
	"os"

	"github.com/ankeesler/btool/app/collector/cc"
	"github.com/ankeesler/btool/log"
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
	root := flag.String("root", ".", "Root path in which to generate")
	project := flag.String("project", "btool", "Name of C++ project")
	path := flag.String("path", "", "Path to generate (no extension)")
	help := flag.Bool("help", false, "Show this help message")

	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(1)
	}

	fs := afero.NewOsFs()
	if err := cc.GenerateClass(fs, *root, *project, *path); err != nil {
		return errors.Wrap(err, "generate class")
	}

	return nil
}
