package resolver_test

import (
	"testing"

	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/resolver"
	"github.com/ankeesler/btool/node/testutil"
	"github.com/go-test/deep"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

var (
	googletestNode = node.Node{
		Name: "googletest",
		Sources: []string{
			"/cache/googletest/googletest-release-1.8.1/googletest/src/gtest-all.cc",
			"/cache/googletest/googletest-release-1.8.1/googlemock/src/gmock-all.cc",
			"/cache/googletest/googletest-release-1.8.1/googlemock/src/gmock_main.cc",
		},
		Headers: []string{
			"gtest/gtest.h",
			"gmock/gmock.h",
		},
		IncludePaths: []string{
			"/cache/googletest/googletest-release-1.8.1/googletest/include",
			"/cache/googletest/googletest-release-1.8.1/googlemock/include",
		},
	}

	googletestNodesWithoutRemoteDependencies = []*node.Node{}

	googletestNodesWithRemoteDependencies = []*node.Node{
		&node.Node{
			Name:         "dep-0/dep-0-test.c",
			Sources:      []string{"dep-0/dep-0-test.c"},
			Headers:      []string{},
			Dependencies: []*node.Node{&testutil.Dep0h, &googletestNode},
		},
		&googletestNode,
	}
)

// TODO: this test won't run when there is no network! fix that!
func TestRemoteHandle(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

	fs := afero.NewMemMapFs()
	if err := afero.WriteFile(
		fs,
		"/dep-0/dep-0-test.c",
		[]byte(`// dep-0-test.c

#include "dep-0/dep-0.h"

#include "gtest/gtest.h"
#include "gmock/gmock.h"
`),
		0644,
	); err != nil {
		t.Fatal(err)
	}

	nodes := []*node.Node{
		&node.Node{
			Name:         "dep-0/dep-0-test.c",
			Sources:      []string{"dep-0/dep-0-test.c"},
			Headers:      []string{},
			Dependencies: []*node.Node{&testutil.Dep0h},
		},
	}

	googletestNode := node.Node{
		Name: "googletest",
		Sources: []string{
			"/cache/googletest/googletest-release-1.8.1/googletest/src/gtest-all.cc",
			"/cache/googletest/googletest-release-1.8.1/googlemock/src/gmock-all.cc",
			"/cache/googletest/googletest-release-1.8.1/googlemock/src/gmock_main.cc",
		},
		Headers: []string{
			"/cache/googletest/googletest-release-1.8.1/googletest/include/gtest/gtest.h",
			"/cache/googletest/googletest-release-1.8.1/googlemock/include/gmock/gmock.h",
		},
		IncludePaths: []string{
			"/cache/googletest/googletest-release-1.8.1/googletest/include",
			"/cache/googletest/googletest-release-1.8.1/googlemock/include",
		},
	}

	exNodes := []*node.Node{
		&node.Node{
			Name:         "dep-0/dep-0-test.c",
			Sources:      []string{"dep-0/dep-0-test.c"},
			Headers:      []string{},
			Dependencies: []*node.Node{&testutil.Dep0h, &googletestNode},
		},
		&googletestNode,
	}

	r := resolver.NewRemote(fs, "/", "/cache")

	acNodes, err := r.Handle(nodes)
	if err != nil {
		t.Fatal(err)
	}

	if diff := deep.Equal(exNodes, acNodes); diff != nil {
		t.Error(diff)
	}
}
