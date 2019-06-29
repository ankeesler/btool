package compiler_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/compiler"
	"github.com/ankeesler/btool/node/compiler/compilerfakes"
	"github.com/ankeesler/btool/node/testutil"
	"github.com/go-test/deep"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func TestHandle(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

	fs := afero.NewMemMapFs()
	c := wireFakeCompiler(fs)
	compiler := compiler.New(c, fs, "/", "/cache")

	nodes := initObjects(testutil.BasicNodes)
	testutil.PopulateFS(nodes, fs)
	exNodes := addObjects(nodes)

	// First build is successful and should load everything into the build cache.
	acNodes, err := compiler.Handle(nodes)
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(exNodes, acNodes); diff != nil {
		t.Error(diff)
	}
	if ex, ac := 3, compileCallCount(c, false); ex != ac {
		t.Errorf("expected %d, got %d", ex, ac)
	}

	// Second build should involve nothing getting re-compiled.
	acNodes, err = compiler.Handle(nodes)
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(exNodes, acNodes); diff != nil {
		t.Error(diff)
	}
	if ex, ac := 3, compileCallCount(c, false); ex != ac {
		t.Errorf("expected %d, got %d", ex, ac)
	}

	// Change the dep-1/dep-1.h header so that main and dep-1 need to be re-compiled.
	if err := afero.WriteFile(
		fs,
		filepath.Join("/", "dep-1/dep-1.h"),
		[]byte("// new data"),
		0600,
	); err != nil {
		t.Fatal(err)
	}

	// Third build should involve main and dep-1 getting re-compiled.
	acNodes, err = compiler.Handle(nodes)
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(exNodes, acNodes); diff != nil {
		t.Error(diff)
	}
	if ex, ac := 5, compileCallCount(c, false); ex != ac {
		t.Errorf("expected %d, got %d", ex, ac)
	}
}

func initObjects(nodes []*node.Node) []*node.Node {
	// TODO: this should be initialized in the testutil package.
	for _, n := range nodes {
		n.Objects = make([]string, 0)
	}
	return nodes
}

func addObjects(nodes []*node.Node) []*node.Node {
	newNodes := testutil.DeepCopy(nodes)
	for _, n := range newNodes {
		for _, s := range n.Sources {
			object := filepath.Join(
				"/cache",
				"objects",
				filepath.Dir(s),
				filepath.Base(s),
			)
			object = object + ".o"
			n.Objects = append(n.Objects, object)
		}
	}
	return newNodes
}

func wireFakeCompiler(fs afero.Fs) *compilerfakes.FakeC {
	c := &compilerfakes.FakeC{}

	c.CompileCStub = func(output, input string, includeDirs []string) error {
		logrus.Debugf("fake c compile %s", output)

		// Sleeping to simulate toolchain calls.
		time.Sleep(time.Millisecond * 100)
		if err := afero.WriteFile(fs, output, []byte("compile c"), 0600); err != nil {
			return errors.Wrap(err, "compile")
		}

		return nil
	}

	c.CompileCCStub = func(output, input string, includeDirs []string) error {
		logrus.Debugf("fake cc compile %s", output)

		// Sleeping to simulate toolchain calls.
		time.Sleep(time.Millisecond * 100)
		if err := afero.WriteFile(fs, output, []byte("compile cc"), 0600); err != nil {
			return errors.Wrap(err, "compile")
		}

		return nil
	}

	return c
}

func compileCallCount(c *compilerfakes.FakeC, cc bool) int {
	if cc {
		return c.CompileCCCallCount()
	} else {
		return c.CompileCCallCount()
	}
}
