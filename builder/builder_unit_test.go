package builder_test

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ankeesler/btool/builder"
	"github.com/ankeesler/btool/config"
	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/testutil"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type fakeCompilerCompileCall struct {
	output, input, root string
}

type fakeCompiler struct {
	fs      afero.Fs
	callsC  []*fakeCompilerCompileCall
	callsCC []*fakeCompilerCompileCall
}

func newFakeCompiler(fs afero.Fs) *fakeCompiler {
	return &fakeCompiler{
		fs:      fs,
		callsC:  make([]*fakeCompilerCompileCall, 0),
		callsCC: make([]*fakeCompilerCompileCall, 0),
	}
}

func (fc *fakeCompiler) CompileC(output, input, root string) error {
	logrus.Debugf("fake c compile %s", output)

	// Sleeping to simulate compiler calls.
	time.Sleep(time.Millisecond * 100)
	if err := afero.WriteFile(fc.fs, output, []byte("compile"), 0600); err != nil {
		return errors.Wrap(err, "link")
	}

	fc.callsC = append(fc.callsC, &fakeCompilerCompileCall{
		output: output,
		input:  input,
		root:   root,
	})
	return nil
}

func (fc *fakeCompiler) CompileCC(output, input, root string) error {
	logrus.Debugf("fake cc compile %s", output)

	// Sleeping to simulate compiler calls.
	time.Sleep(time.Millisecond * 100)
	if err := afero.WriteFile(fc.fs, output, []byte("compile"), 0600); err != nil {
		return errors.Wrap(err, "link")
	}

	fc.callsCC = append(fc.callsCC, &fakeCompilerCompileCall{
		output: output,
		input:  input,
		root:   root,
	})
	return nil
}

func (fc *fakeCompiler) calls(cc bool) []*fakeCompilerCompileCall {
	if cc {
		return fc.callsCC
	} else {
		return fc.callsC
	}
}

type fakeLinkerLinkCall struct {
	output string
	inputs []string
}

type fakeLinker struct {
	fs    afero.Fs
	calls []*fakeLinkerLinkCall
}

func newFakeLinker(fs afero.Fs) *fakeLinker {
	return &fakeLinker{
		fs:    fs,
		calls: make([]*fakeLinkerLinkCall, 0),
	}
}

func (fl *fakeLinker) Link(output string, inputs []string) error {
	logrus.Debugf("fake link %s", output)

	// Sleeping to simulate linker calls.
	time.Sleep(time.Millisecond * 100)
	if err := afero.WriteFile(fl.fs, output, []byte("link"), 0700); err != nil {
		return errors.Wrap(err, "link")
	}

	fl.calls = append(fl.calls, &fakeLinkerLinkCall{
		output: output,
		inputs: inputs,
	})
	return nil
}

func TestBuildUnit(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

	projects := []*testutil.Project{
		testutil.BasicProjectC(),
		testutil.BasicProjectCC(),
	}

	for _, project := range projects {
		t.Run(project.Name, func(t *testing.T) {
			fs := afero.NewMemMapFs()

			project.Root = "/tmp/root"
			if err := project.PopulateFS(fs); err != nil {
				t.Fatal(err)
			}

			cc := strings.HasSuffix(project.Name, "CC")

			cfg := config.Config{
				Name:  "some-project-name",
				Root:  project.Root,
				Cache: "/some/cache/root",
			}
			c := newFakeCompiler(fs)
			l := newFakeLinker(fs)
			b := builder.New(fs, &cfg, c, l)

			// First build is successful and should load everything into the build cache.
			if err := b.Build(project.Graph()); err != nil {
				t.Fatal(err)
			}
			if ex, ac := 2, len(c.calls(cc)); ex != ac {
				t.Errorf("expected %d, got %d", ex, ac)
			}
			if ex, ac := 1, len(l.calls); ex != ac {
				t.Errorf("expected %d, got %d", ex, ac)
			}

			// Second build should involve nothing getting re-compiled.
			if err := b.Build(project.Graph()); err != nil {
				t.Fatal(err)
			}
			if ex, ac := 2, len(c.calls(cc)); ex != ac {
				t.Errorf("expected %d, got %d", ex, ac)
			}
			if ex, ac := 1, len(l.calls); ex != ac {
				t.Errorf("expected %d, got %d", ex, ac)
			}

			// Change the master.h header so that main.cc needs to be re-compiled.
			if err := afero.WriteFile(
				fs,
				filepath.Join(project.Root, "master.h"),
				[]byte("// new data"),
				0600,
			); err != nil {
				t.Fatal(err)
			}

			// Third build should involve main.cc getting re-compiled, and stuff to
			// get re-linked.
			if err := b.Build(project.Graph()); err != nil {
				t.Fatal(err)
			}
			if ex, ac := 3, len(c.calls(cc)); ex != ac {
				t.Errorf("expected %d, got %d", ex, ac)
			}
			if ex, ac := 2, len(l.calls); ex != ac {
				t.Errorf("expected %d, got %d", ex, ac)
			}
		})
	}
}
