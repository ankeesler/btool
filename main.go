package main

import (
	command "github.com/ankeesler/btool/cmd"
	"github.com/ankeesler/btool/formatter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(formatter.New())
	if err := run(); err != nil {
		logrus.Fatal(err)
	}
}

func run() error {
	cmd, err := command.Init()
	if err != nil {
		return errors.Wrap(err, "init")
	}

	if err := cmd.Execute(); err != nil {
		return errors.Wrap(err, "run")
	}

	return nil
}
