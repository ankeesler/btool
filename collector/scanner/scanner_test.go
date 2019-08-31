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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScannerCollect(t *testing.T) {
	fs := afero.NewMemMapFs()
	root := "/some/root"
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

	data := []struct {
		name           string
		target         string
		exNFunc        func(*collector.Ctx) *node.Node
		exIncludePaths [][]string
	}{
		{
			name:   "BasicC",
			target: "/some/root/main",
			exNFunc: func(ctx *collector.Ctx) *node.Node {
				exDep0HN := node.New("/some/root/dep-0/dep-0.h")
				exDep1HN := node.New("/some/root/dep-1/dep-1.h").Dependency(exDep0HN)
				exDep0CN := node.New("/some/root/dep-0/dep-0.c").Dependency(exDep0HN)
				exOutsideDepNH := node.New("/some/other/root/outside/dep.h")
				exDep1CN := node.New("/some/root/dep-1/dep-1.c").Dependency(
					exDep1HN, exDep0HN, exOutsideDepNH)
				exMainCN := node.New("/some/root/main.c").Dependency(exDep1HN, exDep0HN)
				exDep0ON := node.New("/some/root/dep-0/dep-0.o").Dependency(exDep0CN)
				exDep0ON.Resolver = compileCR
				exDep1ON := node.New("/some/root/dep-1/dep-1.o").Dependency(exDep1CN)
				exDep1ON.Resolver = compileCR
				exMainON := node.New("/some/root/main.o").Dependency(exMainCN)
				exMainON.Resolver = compileCR
				exMainN := node.New("/some/root/main").Dependency(exMainON, exDep1ON, exDep0ON)
				exMainN.Resolver = linkCR

				ctx.AddIncludePath("/some/other/root/outside")
				ctx.NS.Add(exOutsideDepNH)

				return exMainN
			},
			exIncludePaths: [][]string{
				[]string{
					root,
				},
				[]string{
					root,
					"/some/other/root/outside/",
				},
				[]string{
					root,
				},
			},
		},
		{
			name:   "Test",
			target: "a_test",
			exNFunc: func(ctx *collector.Ctx) *node.Node {
				gtestH := node.New(".btool/abc123/gtest.h")
				aHN := node.New("a.h")
				aCCN := node.New("a.cc").Dependency(aHN)
				aON := node.New("a.o").Dependency(aCCN)
				aON.Resolver = compileCCR
				aTestCCN := node.New("a_test.cc").Dependency(aHN, gtestH)
				aTestON := node.New("a_test.o").Dependency(aTestCCN)
				aTestON.Resolver = compileCCR
				aTestN := node.New("a_test").Dependency(aTestON, aON)
				aTestN.Resolver = linkCCR

				ctx.AddIncludePath(".btool/abc123")
				ctx.NS.Add(gtestH)

				return aTestN
			},
			exIncludePaths: [][]string{
				[]string{ // a.cc
					root,
				},
				[]string{ // a_test.cc
					root,
					".btool/abc123",
				},
			},
		},
	}
	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			ns := collector.NewNodeStore(nil)
			ctx := collector.NewCtx(ns, rf)

			exN := datum.exNFunc(ctx)
			populateFS(t, exN, fs, root)
			printFS(t, fs, root)

			s := scanner.New(fs, root, i)

			acN := node.New(datum.target)
			require.Nil(t, s.Collect(ctx, acN))

			sorter := sorter.New()
			sorter.Collect(&collector.Ctx{}, acN)
			sorter.Collect(&collector.Ctx{}, exN)
			log.Debugf("%s\n%s", "ac", node.String(acN))
			log.Debugf("%s\n%s", "ex", node.String(exN))
			require.Nil(t, deep.Equal(exN, acN))

			require.Equal(t, len(datum.exIncludePaths), rf.NewCompileCCallCount())
			for i, exIncludePaths := range datum.exIncludePaths {
				assert.Equal(
					t,
					exIncludePaths,
					rf.NewCompileCArgsForCall(i),
				)
			}
		})
	}
}

func populateFS(t *testing.T, n *node.Node, fs afero.Fs, root string) {
	require.Nil(
		t,
		node.Visit(
			n,
			func(nN *node.Node) error {
				return reallyPopulateFS(nN, fs, root)
			},
		),
	)
}

func printFS(t *testing.T, fs afero.Fs, root string) {
	require.Nil(
		t,
		afero.Walk(
			fs,
			root,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

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
}

func reallyPopulateFS(n *node.Node, fs afero.Fs, root string) error {
	if !strings.HasSuffix(n.Name, ".c") && !strings.HasSuffix(n.Name, ".h") {
		return nil
	}

	b := bytes.NewBuffer([]byte{})
	fmt.Fprintf(b, "// %s\n", n.Name)

	for _, dN := range n.Dependencies {
		if strings.HasSuffix(dN.Name, ".h") {
			var header string
			if strings.HasPrefix(dN.Name, root) {
				var err error
				header, err = filepath.Rel(root, dN.Name)
				if err != nil {
					return errors.Wrap(err, "rel")
				}
			} else {
				header = filepath.Base(dN.Name)
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
