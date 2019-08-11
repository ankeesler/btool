// Package testutil provides test utilities for the registry package.
package testutil

import "github.com/ankeesler/btool/node/registry"

func Index() *registry.Index {
	return &registry.Index{
		Files: []registry.IndexFile{
			registry.IndexFile{
				Path:   "file_btool.yml",
				SHA256: "92c65ba8e54870136bc34b1aae13a69932a8fbc1b89bfdd04751be00f7d13352",
			},
			registry.IndexFile{
				Path:   "some/path/to/file_btool.yml",
				SHA256: "92c65ba8e54870136bc34b1aae13a69932a8fbc1b89bfdd04751be00f7d13352",
			},
		},
	}
}

func Nodes() []*registry.Node {
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
