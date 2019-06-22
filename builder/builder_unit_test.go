package builder_test

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ankeesler/btool/builder"
	"github.com/ankeesler/btool/builder/builderfakes"
	"github.com/ankeesler/btool/config"
	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/testutil"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

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
			tc := wireFakeToolchain(fs)
			b := builder.New(fs, &cfg, tc)

			// First build is successful and should load everything into the build cache.
			if err := b.Build(project.Graph()); err != nil {
				t.Fatal(err)
			}
			if ex, ac := 2, compileCallCount(tc, cc); ex != ac {
				t.Errorf("expected %d, got %d", ex, ac)
			}
			if ex, ac := 1, tc.LinkCallCount(); ex != ac {
				t.Errorf("expected %d, got %d", ex, ac)
			}

			// Second build should involve nothing getting re-compiled.
			if err := b.Build(project.Graph()); err != nil {
				t.Fatal(err)
			}
			if ex, ac := 2, compileCallCount(tc, cc); ex != ac {
				t.Errorf("expected %d, got %d", ex, ac)
			}
			if ex, ac := 1, tc.LinkCallCount(); ex != ac {
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
			if ex, ac := 3, compileCallCount(tc, cc); ex != ac {
				t.Errorf("expected %d, got %d", ex, ac)
			}
			if ex, ac := 2, tc.LinkCallCount(); ex != ac {
				t.Errorf("expected %d, got %d", ex, ac)
			}
		})
	}
}

func wireFakeToolchain(fs afero.Fs) *builderfakes.FakeToolchain {
	t := &builderfakes.FakeToolchain{}

	t.CompileCStub = func(output, input string, includeDirs []string) error {
		logrus.Debugf("fake c compile %s", output)

		// Sleeping to simulate toolchain calls.
		time.Sleep(time.Millisecond * 100)
		if err := afero.WriteFile(fs, output, []byte("compile"), 0600); err != nil {
			return errors.Wrap(err, "link")
		}

		return nil
	}

	t.CompileCCStub = func(output, input string, includeDirs []string) error {
		logrus.Debugf("fake c compile %s", output)

		// Sleeping to simulate toolchain calls.
		time.Sleep(time.Millisecond * 100)
		if err := afero.WriteFile(fs, output, []byte("compile"), 0600); err != nil {
			return errors.Wrap(err, "compile")
		}

		return nil
	}

	t.LinkStub = func(output string, inputs []string) error {
		logrus.Debugf("fake link %s", output)

		// Sleeping to simulate toolchain calls.
		time.Sleep(time.Millisecond * 100)
		if err := afero.WriteFile(fs, output, []byte("link"), 0700); err != nil {
			return errors.Wrap(err, "link")
		}

		return nil
	}

	return t
}

func compileCallCount(t *builderfakes.FakeToolchain, cc bool) int {
	if cc {
		return t.CompileCCCallCount()
	} else {
		return t.CompileCCallCount()
	}
}
