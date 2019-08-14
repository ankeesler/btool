// Package testutil provides test utilities for the registry package.
package testutil

import "github.com/ankeesler/btool/node/registry"

func Index() *registry.Index {
	return &registry.Index{
		Name: "some-index",
		Files: []registry.IndexFile{
			registry.IndexFile{
				Path:   "file_a_btool.yml",
				SHA256: "3f67ca5a3aa3f9673a8e76cca937b4b9e1de527ed23d6d12f986c4d63b2c88a2",
			},
			registry.IndexFile{
				Path:   "some/path/to/file_b_btool.yml",
				SHA256: "f31d694f3ee869d87126a7e472c0815ad81b5994a86bd9b19a31d32c5a83c485",
			},
		},
	}
}

func FileAGaggle() *registry.Gaggle {
	return &registry.Gaggle{
		Metadata: map[string]interface{}{
			"foo":     "bar",
			"project": "some-project",
		},
		Nodes: []*registry.Node{
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
		},
	}
}

func FileBGaggle() *registry.Gaggle {
	return &registry.Gaggle{
		Metadata: map[string]interface{}{
			"bat":     "mayonnaise",
			"project": "some-other-project",
		},
		Nodes: []*registry.Node{
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
		},
	}
}
