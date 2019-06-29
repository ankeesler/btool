package testutil

import "github.com/ankeesler/btool/node"

var (
	Dep0h = node.Node{
		Name:         "dep-0/dep-0.h",
		Sources:      []string{},
		Headers:      []string{"dep-0/dep-0.h"},
		Dependencies: []*node.Node{},
	}
	Dep0c = node.Node{
		Name:         "dep-0/dep-0.c",
		Sources:      []string{"dep-0/dep-0.c"},
		Headers:      []string{},
		Dependencies: []*node.Node{&Dep0h},
	}

	Dep1h = node.Node{
		Name:         "dep-1/dep-1.h",
		Sources:      []string{},
		Headers:      []string{"dep-1/dep-1.h"},
		Dependencies: []*node.Node{&Dep0h},
	}
	Dep1c = node.Node{
		Name:         "dep-1/dep-1.c",
		Sources:      []string{"dep-1/dep-1.c"},
		Headers:      []string{},
		Dependencies: []*node.Node{&Dep1h, &Dep0h},
	}

	Mainc = node.Node{
		Name:         "main.c",
		Sources:      []string{"main.c"},
		Headers:      []string{},
		Dependencies: []*node.Node{&Dep1h, &Dep0h},
	}
)

var (
	BasicNodes = []*node.Node{
		&Dep0c,
		&Dep0h,
		&Dep1c,
		&Dep1h,
		&Mainc,
	}
)
