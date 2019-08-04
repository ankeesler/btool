// Package main provides a main function that reads in a bunch of nodes and
// applies a resolver.NoopResolver to show the resolution.
package main

import (
	"flag"
	"os"

	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/builder"
	"github.com/ankeesler/btool/node/deps"
	"github.com/ankeesler/btool/node/objecter"
	"github.com/ankeesler/btool/node/sorter"
	"github.com/ankeesler/btool/node/walker"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func main() {
	root := flag.String("root", "/tmp/BasicC", "Root of node list")
	cache := flag.String("cache", "/tmp/btool-node-test", "Cache directory")
	target := flag.String("target", "main", "Target to build")
	help := flag.Bool("help", false, "Show this help message")

	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(1)
	}

	cfg := node.Config{
		Root:   *root,
		Cache:  *cache,
		Target: *target,

		CCompiler:  "clang",
		CCCompiler: "clang++",
	}

	fs := afero.NewOsFs()

	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

	if err := node.Pipeline(
		&cfg,
		walker.New(fs),
		deps.NewLocal(fs),
		objecter.New(),
		sorter.NewAlpha(),
		builder.New(),
	); err != nil {
		logrus.Error(err)
	}
}
