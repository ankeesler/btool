// This program provides a regeneration mechanism for test fixtures.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ankeesler/btool/node/testutil"
	"github.com/spf13/afero"
)

func main() {
	root := flag.String("root", "/tmp", "Root at which to generate the fixture")
	help := flag.Bool("help", false, "Show this help message")

	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(1)
	}

	fs := afero.NewOsFs()

	fixtures := map[string]testutil.Nodes{
		"BasicC":  testutil.BasicNodesC,
		"BasicCC": testutil.BasicNodesCC,
	}

	for fixture, nodes := range fixtures {
		fixtureRoot := filepath.Join(*root, fixture)
		if err := os.RemoveAll(fixtureRoot); err != nil {
			fmt.Printf("error: %s", err.Error())
			os.Exit(1)
		}
		fmt.Printf("Generating %s...\n", fixtureRoot)

		nodes.PopulateFS(fixtureRoot, fs)
	}
}
