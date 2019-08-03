package testutil

import "github.com/ankeesler/btool/node"

var (
	Dep0h  = node.New("dep-0/dep-0.h")
	Dep0c  = node.New("dep-0/dep-0.c").Dependency(Dep0h)
	Dep0cc = node.New("dep-0/dep-0.cc").Dependency(Dep0h)

	Dep0co  = node.New("dep-0/dep-0.o").Dependency(Dep0c)
	Dep0cco = node.New("dep-0/dep-0.o").Dependency(Dep0cc)

	Dep1h  = node.New("dep-1/dep-1.h").Dependency(Dep0h)
	Dep1c  = node.New("dep-1/dep-1.c").Dependency(Dep1h, Dep0h)
	Dep1cc = node.New("dep-1/dep-1.cc").Dependency(Dep1h, Dep0h)

	Dep1co  = node.New("dep-1/dep-1.o").Dependency(Dep1c)
	Dep1cco = node.New("dep-1/dep-1.o").Dependency(Dep1cc)

	Mainc  = node.New("main.c").Dependency(Dep1h, Dep0h)
	Maincc = node.New("main.cc").Dependency(Dep1h, Dep0h)

	Mainco  = node.New("main.o").Dependency(Mainc)
	Maincco = node.New("main.o").Dependency(Maincc)
)

var (
	BasicNodesC Nodes = []*node.Node{
		Dep0c,
		Dep0h,
		Dep1c,
		Dep1h,
		Mainc,
	}

	BasicNodesCO Nodes = []*node.Node{
		Dep0c,
		Dep0h,
		Dep1c,
		Dep1h,
		Mainc,
		Dep0co,
		Dep1co,
		Mainco,
	}

	BasicNodesCC Nodes = []*node.Node{
		Dep0cc,
		Dep0h,
		Dep1cc,
		Dep1h,
		Maincc,
	}

	BasicNodesCCO Nodes = []*node.Node{
		Dep0cc,
		Dep0h,
		Dep1cc,
		Dep1h,
		Maincc,
		Dep0cco,
		Dep1cco,
		Maincco,
	}
)
