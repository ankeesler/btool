package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ankeesler/btool/builder"
	"github.com/ankeesler/btool/builder/toolchain"
	"github.com/ankeesler/btool/config"
	"github.com/ankeesler/btool/deps"
	"github.com/ankeesler/btool/deps/downloader"
	"github.com/ankeesler/btool/generator"
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
		cache    string
		logLevel string

		s *scanner.Scanner
		b *builder.Builder
		g *generator.Generator
	)

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

			fs := afero.NewOsFs()
			c := config.Config{
				Name:  filepath.Base(root),
				Root:  root,
				Cache: cache,
			}
			downloader := downloader.New(func(file string) bool {
				return strings.HasSuffix(file, ".c") ||
					strings.HasSuffix(file, ".cc") ||
					strings.HasSuffix(file, ".h")
			})
			d := deps.New(fs, cache, downloader)

			s = scanner.New(fs, &c, d)
			b = builder.New(
				fs,
				&c,
				toolchain.New("clang", "clang++", "clang"),
			)
			g = generator.New(fs, &c)

			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			logrus.Infof("success")
			return nil
		},
	}
	rootCmdFlags := rootCmd.PersistentFlags()
	rootCmdFlags.StringVar(&root, "root", filepath.Join(cwd, ".btool"), "Path to project root")
	rootCmdFlags.StringVar(&cache, "cache", os.TempDir(), "Path to build cache")
	rootCmdFlags.StringVar(&logLevel, "loglevel", "info", "Log level")

	scanCmd := &cobra.Command{
		Use:   "scan [file]",
		Short: "Display build graph",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var g *graph.Graph
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
			g, err := s.ScanFile(args[0])
			if err != nil {
				return errors.Wrap(err, "scan")
			}

			if err := b.Build(g); err != nil {
				return errors.Wrap(err, "build")
			}

			return nil
		},
	}
	rootCmd.AddCommand(buildCmd)

	cleanCmd := &cobra.Command{
		Use:   "clean",
		Short: "Clean build data",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := b.Clean(); err != nil {
				return errors.Wrap(err, "clean")
			}

			return nil
		},
	}
	rootCmd.AddCommand(cleanCmd)

	classCmd := &cobra.Command{
		Use:     "class <path/without/file/extension>",
		Short:   "Generate class",
		Example: "btool class path/to/some/class/without/extension",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := g.Class(args[0]); err != nil {
				return errors.Wrap(err, "class")
			}

			return nil
		},
	}
	rootCmd.AddCommand(classCmd)

	return rootCmd, nil
}
