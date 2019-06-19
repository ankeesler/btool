package deps_test

import (
	"strings"
	"testing"

	"github.com/ankeesler/btool/deps"
	"github.com/ankeesler/btool/deps/downloader"
	"github.com/ankeesler/btool/formatter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func TestResolveInclude(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

	fs := afero.NewMemMapFs()
	cache := "/tuna/cache"
	downloader := downloader.New(func(s string) bool {
		return strings.HasSuffix(s, ".h")
	})
	d := deps.New(fs, cache, downloader)

	// unknown include
	if include, err := d.ResolveInclude("this/does/not/exist.h"); err != nil {
		t.Error(err)
	} else if include != "" {
		t.Errorf("expected '', got '%s'", include)
	}

	// googletest
	include, err := d.ResolveInclude("gtest/gtest.h")
	if err != nil {
		t.Error(err)
	} else if exists, err := afero.Exists(fs, include); err != nil {
		t.Error(err)
	} else if !exists {
		t.Error("gtest: expected file to exist:", include)
	}

	include, err = d.ResolveInclude("gmock/gmock.h")
	if err != nil {
		t.Error(err)
	} else if exists, err := afero.Exists(fs, include); err != nil {
		t.Error(err)
	} else if !exists {
		t.Error("gmock: expected file to exist:", include)
	}
}

func TestResolveSources(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

	fs := afero.NewMemMapFs()
	cache := "/tuna/cache"
	downloader := downloader.New(func(s string) bool {
		return strings.HasSuffix(s, ".cc")
	})
	d := deps.New(fs, cache, downloader)

	// unknown include
	if sources, err := d.ResolveSources("this/does/not/exist.h"); err != nil {
		t.Error(err)
	} else if sources != nil {
		t.Error("expected <nil>, got", sources)
	}

	// googletest
	sources, err := d.ResolveSources("gtest/gtest.h")
	if err != nil {
		t.Error(err)
	} else {
		for _, source := range sources {
			if exists, err := afero.Exists(fs, source); err != nil {
				t.Error(err)
			} else if !exists {
				t.Error("gtest: expected file to exist:", source)
			}
		}
	}

	sources, err = d.ResolveSources("gmock/gmock.h")
	if err != nil {
		t.Error(err)
	} else {
		for _, source := range sources {
			if exists, err := afero.Exists(fs, source); err != nil {
				t.Error(err)
			} else if !exists {
				t.Error("gmock: expected file to exist:", source)
			}
		}
	}
}
