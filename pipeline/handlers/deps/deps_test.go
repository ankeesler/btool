package deps_test

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ankeesler/btool/formatter"
	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/resolvers"
	"github.com/ankeesler/btool/node/testutil"
	"github.com/ankeesler/btool/pipeline"
	"github.com/ankeesler/btool/pipeline/handlers/deps"
	"github.com/go-test/deep"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func TestDeps(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter.New())

	data, err := ioutil.ReadFile(filepath.Join("testdata", "btool_deps.yml"))
	if err != nil {
		t.Fatal(err)
	}

	var v deps.Version
	if err := yaml.Unmarshal(data, &v); err != nil {
		t.Fatal(err)
	}

	if ex, ac := "0.0.1", v.Version; ex != ac {
		t.Error(ex, "!=", ac)
	}

	var d deps.Deps
	if err := yaml.Unmarshal(data, &d); err != nil {
		t.Fatal(err)
	}
	t.Log("d", d)

	ctx := pipeline.NewCtxBuilder().Build()
	d.Nodes(ctx)
	if ctx.Err != nil {
		t.Fatal(ctx.Err)
	}
	acNodes := ctx.Nodes

	exNodes := testutil.BasicNodesC.Copy()
	exNodes, err = addObjects(exNodes)
	if err != nil {
		t.Fatal(err)
	}
	exNodes, err = addLibrary(exNodes)
	if err != nil {
		t.Fatal(err)
	}
	exNodes, err = addExecutable(exNodes)
	if err != nil {
		t.Fatal(err)
	}

	node.SortAlpha(exNodes)
	node.SortAlpha(acNodes)
	if diff := deep.Equal(exNodes.Cast(), acNodes); diff != nil {
		t.Error(diff)
	}
}

func addObjects(nodes []*node.Node) ([]*node.Node, error) {
	objects := []string{"main.o", "dep-0/dep-0.o", "dep-1/dep-1.o"}
	for _, o := range objects {
		c := strings.ReplaceAll(o, ".o", ".c")
		cN := node.Find(c, nodes)
		if cN == nil {
			return nil, errors.New(c + " is nil")
		}

		oN := node.New(o).Dependency(cN)
		oN.Resolver = resolvers.NewCompile("", "", nil)
		nodes = append(nodes, oN)
	}
	return nodes, nil
}

func addLibrary(nodes []*node.Node) ([]*node.Node, error) {
	dep0oN := node.Find("dep-0/dep-0.o", nodes)
	if dep0oN == nil {
		return nil, errors.New("dep0oN is nil")
	}

	dep1oN := node.Find("dep-1/dep-1.o", nodes)
	if dep1oN == nil {
		return nil, errors.New("dep1oN is nil")
	}

	l := "dep.a"
	lN := node.New(l).Dependency(dep0oN, dep1oN)
	lN.Resolver = resolvers.NewArchive("", "")
	nodes = append(nodes, lN)

	return nodes, nil
}

func addExecutable(nodes []*node.Node) ([]*node.Node, error) {
	mainoN := node.Find("main.o", nodes)
	if mainoN == nil {
		return nil, errors.New("main.o is nil")
	}

	dep0oN := node.Find("dep-0/dep-0.o", nodes)
	if dep0oN == nil {
		return nil, errors.New("dep-0/dep-0.o is nil")
	}

	dep1oN := node.Find("dep-1/dep-1.o", nodes)
	if dep1oN == nil {
		return nil, errors.New("dep-1/dep-1.o is nil")
	}

	mainN := node.New("main").Dependency(mainoN, dep0oN, dep1oN)
	mainN.Resolver = resolvers.NewLink("", "")
	nodes = append(nodes, mainN)

	return nodes, nil
}
