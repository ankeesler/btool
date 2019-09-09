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
		root           string
		exNFunc        func(*collector.Ctx) *node.Node
		exIncludePaths [][]string
		cc             bool
	}{
		{
			name:   "BasicC",
			target: "/some/root/main",
			root:   "/some/root",
			exNFunc: func(ctx *collector.Ctx) *node.Node {
				exDep0HN := node.New("/some/root/dep-0/dep-0.h")
				exDep1HN := node.New("/some/root/dep-1/dep-1.h").Dependency(exDep0HN)
				exDep0CN := node.New("/some/root/dep-0/dep-0.c").Dependency(exDep0HN)
				exOutsideDepNH := node.New(".btool/abc123/some/other/root/outside/dep.h")
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

				ctx.AddIncludePath(".btool/abc123/some/other/root/outside")
				ctx.NS.Add(exOutsideDepNH)

				return exMainN
			},
			exIncludePaths: [][]string{
				[]string{ // dep-1/dep-1.c
					"/some/root",
				},
				[]string{ // dep-0/dep-0.c
					"/some/root",
					".btool/abc123/some/other/root/outside/",
				},
				[]string{ // main.c
					"/some/root",
				},
			},
			cc: false,
		},
		{
			name:   "Test",
			target: "a_test",
			root:   "",
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
					"",
				},
				[]string{ // a_test.cc
					".btool/abc123/",
				},
			},
			cc: true,
		},
	}
	for _, datum := range data {
		t.Run(datum.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			i := includeser.New(fs)

			ns := collector.NewNodeStore(nil)
			ctx := collector.NewCtx(ns, rf)

			exN := datum.exNFunc(ctx)
			populateFS(t, exN, fs, datum.root)
			printFS(t, fs, datum.root)

			s := scanner.New(fs, datum.root, i)

			acN := node.New(datum.target)
			require.Nil(t, s.Collect(ctx, acN))

			sorter := sorter.New()
			sorter.Collect(&collector.Ctx{}, acN)
			sorter.Collect(&collector.Ctx{}, exN)
			log.Debugf("%s\n%s", "ac", node.String(acN))
			log.Debugf("%s\n%s", "ex", node.String(exN))
			require.Nil(t, deep.Equal(exN, acN))

			var compileCallCountFunc func() int
			var compileArgsForCallFunc func(int) []string
			if datum.cc {
				compileCallCountFunc = rf.NewCompileCCCallCount
				compileArgsForCallFunc = rf.NewCompileCCArgsForCall
			} else {
				compileCallCountFunc = rf.NewCompileCCallCount
				compileArgsForCallFunc = rf.NewCompileCArgsForCall
			}

			require.Equal(t, len(datum.exIncludePaths), compileCallCountFunc())
			for i, exIncludePaths := range datum.exIncludePaths {
				assert.Equal(
					t,
					exIncludePaths,
					compileArgsForCallFunc(i),
					"call #%d",
					i,
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
					return errors.Wrap(err, "err")
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
	if !strings.HasSuffix(n.Name, ".c") &&
		!strings.HasSuffix(n.Name, ".cc") &&
		!strings.HasSuffix(n.Name, ".h") {
		return nil
	}

	b := bytes.NewBuffer([]byte{})
	fmt.Fprintf(b, "// %s\n", n.Name)

	for _, dN := range n.Dependencies {
		if strings.HasSuffix(dN.Name, ".h") {
			var header string
			if !strings.HasPrefix(dN.Name, ".btool") {
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
