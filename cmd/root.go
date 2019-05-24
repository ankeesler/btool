package cmd

import (
	"os"
	"path/filepath"

	"github.com/ankeesler/btool/builder"
	"github.com/ankeesler/btool/builder/compiler"
	"github.com/ankeesler/btool/builder/linker"
	"github.com/ankeesler/btool/scanner"
	"github.com/ankeesler/btool/scanner/graph"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func Init() (*cobra.Command, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "getwd")
	}

	var (
		root     string
		store    string
		logLevel string
	)

	fs := afero.NewOsFs()

	rootCmd := &cobra.Command{
		Use:           "btool",
		Short:         "Push to start C/C++ building!",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			level, err := logrus.ParseLevel(logLevel)
			if err != nil {
				return errors.Wrap(err, "parse log level")
			}
			logrus.SetLevel(level)

			return nil
		},
	}
	rootCmdFlags := rootCmd.PersistentFlags()
	rootCmdFlags.StringVar(&root, "root", cwd, "Path to project root")
	rootCmdFlags.StringVar(&store, "store", filepath.Join(cwd, ".btool"), "Path to btool store")
	rootCmdFlags.StringVar(&logLevel, "loglevel", "info", "Log level")

	scanCmd := &cobra.Command{
		Use:   "scan [file]",
		Short: "Display build graph",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var g *graph.Graph
			s := scanner.New(fs, root)
			if len(args) == 0 {
				g, err = s.ScanRoot()
			} else {
				g, err = s.ScanFile(args[0])
			}
			if err != nil {
				return errors.Wrap(err, "scan")
			}

			logrus.Infof("graph:\n%s", g)

			return nil
		},
	}
	rootCmd.AddCommand(scanCmd)

	buildCmd := &cobra.Command{
		Use:   "build <file>",
		Short: "Build file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s := scanner.New(fs, root)
			g, err := s.ScanFile(args[0])
			if err != nil {
				return errors.Wrap(err, "scan")
			}

			b := builder.New(
				fs,
				root,
				store,
				compiler.New(),
				linker.New(),
			)
			if err := b.Build(g); err != nil {
				return errors.Wrap(err, "build")
			}

			return nil
		},
	}
	rootCmd.AddCommand(buildCmd)

	return rootCmd, nil
}
