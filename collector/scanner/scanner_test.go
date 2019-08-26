package scanner_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ankeesler/btool/collector"
	"github.com/ankeesler/btool/collector/collectorfakes"
	"github.com/ankeesler/btool/collector/scanner"
	"github.com/ankeesler/btool/collector/scanner/includeser"
	"github.com/ankeesler/btool/collector/sorter"
	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/nodefakes"
	"github.com/go-test/deep"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestScannerCollect(t *testing.T) {
	fs := afero.NewMemMapFs()
	root := "/some/root"
	ns := collector.NewNodeStore(nil)
	i := includeser.New(fs)
	linkCR := &nodefakes.FakeResolver{}
	linkCCR := &nodefakes.FakeResolver{}
	compileCR := &nodefakes.FakeResolver{}
	compileCCR := &nodefakes.FakeResolver{}
	rf := &collectorfakes.FakeResolverFactory{}
	rf.NewLinkCReturns(linkCR)
	rf.NewLinkCCReturns(linkCCR)
	rf.NewCompileCReturns(compileCR)
	rf.NewCompileCCReturns(compileCCR)

	exDep0HN := node.New("/some/root/dep-0/dep-0.h")
	exDep1HN := node.New("/some/root/dep-1/dep-1.h").Dependency(exDep0HN)
	exDep0CN := node.New("/some/root/dep-0/dep-0.c").Dependency(exDep0HN)
	exDep1CN := node.New("/some/root/dep-1/dep-1.c").Dependency(exDep1HN, exDep0HN)
	exMainCN := node.New("/some/root/main.c").Dependency(exDep1HN, exDep0HN)
	exDep0ON := node.New("/some/root/dep-0/dep-0.o").Dependency(exDep0CN)
	exDep0ON.Resolver = compileCR
	exDep1ON := node.New("/some/root/dep-1/dep-1.o").Dependency(exDep1CN)
	exDep1ON.Resolver = compileCR
	exMainON := node.New("/some/root/main.o").Dependency(exMainCN)
	exMainON.Resolver = compileCR
	exMainN := node.New("/some/root/main").Dependency(exMainON, exDep1ON, exDep0ON)
	exMainN.Resolver = linkCR
	require.Nil(
		t,
		node.Visit(
			exMainN,
			func(n *node.Node) error {
				return populateFS(n, fs, root)
			},
		),
	)
	require.Nil(
		t,
		afero.Walk(
			fs,
			root,
			func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					content, err := afero.ReadFile(fs, path)
					if err != nil {
						return errors.Wrap(err, "read file")
					}

					log.Debugf("\n%s:\n%s", path, string(content))
				}

				return nil
			},
		),
	)

	s := scanner.New(fs, root, i)

	ctx := collector.NewCtx(ns, rf)
	acMainN := node.New("/some/root/main")
	require.Nil(t, s.Collect(ctx, acMainN))

	sorter := sorter.New()
	sorter.Collect(&collector.Ctx{}, acMainN)
	sorter.Collect(&collector.Ctx{}, exMainN)
	log.Debugf("%s\n%s", "ac", node.String(acMainN))
	log.Debugf("%s\n%s", "ex", node.String(exMainN))
	require.Nil(t, deep.Equal(exMainN, acMainN))
}

func populateFS(n *node.Node, fs afero.Fs, root string) error {
	if !strings.HasSuffix(n.Name, ".c") && !strings.HasSuffix(n.Name, ".h") {
		return nil
	}

	b := bytes.NewBuffer([]byte{})
	fmt.Fprintf(b, "// %s\n", n.Name)

	for _, dN := range n.Dependencies {
		if strings.HasSuffix(dN.Name, ".h") {
			header, err := filepath.Rel(root, dN.Name)
			if err != nil {
				return errors.Wrap(err, "rel")
			}

			if filepath.Dir(dN.Name) == filepath.Dir(n.Name) {
				header = filepath.Base(header)
			}

			fmt.Fprintf(b, "#include \"%s\"\n", header)
		}
	}

	fmt.Fprintf(b, "\n")

	if err := afero.WriteFile(fs, n.Name, b.Bytes(), 0644); err != nil {
		return errors.Wrap(err, "write file")
	}

	return nil
}
