// Package testutil provides test utilities for the registry package.
package testutil

import "github.com/ankeesler/btool/node/registry"

func Index() *registry.Index {
	return &registry.Index{
		Files: []registry.IndexFile{
			registry.IndexFile{
				Path:   "file_a_btool.yml",
				SHA256: "92c65ba8e54870136bc34b1aae13a69932a8fbc1b89bfdd04751be00f7d13352",
			},
			registry.IndexFile{
				Path:   "some/path/to/file_b_btool.yml",
				SHA256: "cfac8628fca088ff0cc5756789ad441c6fd9e5662bb646199f1876207d823873",
			},
		},
	}
}

func FileANodes() []*registry.Node {
	return []*registry.Node{
		&registry.Node{
			Name:         "tuna",
			Dependencies: []string{},
			Resolver: registry.Resolver{
				Name:   "",
				Config: map[string]interface{}{},
			},
		},
		&registry.Node{
			Name:         "fish",
			Dependencies: []string{"tuna"},
			Resolver: registry.Resolver{
				Name:   "",
				Config: map[string]interface{}{},
			},
		},
	}
}

func FileBNodes() []*registry.Node {
	return []*registry.Node{
		&registry.Node{
			Name:         "marlin",
			Dependencies: []string{},
			Resolver: registry.Resolver{
				Name:   "",
				Config: map[string]interface{}{},
			},
		},
		&registry.Node{
			Name:         "bacon",
			Dependencies: []string{"marlin"},
			Resolver: registry.Resolver{
				Name:   "",
				Config: map[string]interface{}{},
			},
		},
	}
}
