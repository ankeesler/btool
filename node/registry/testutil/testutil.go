// Package testutil provides test utilities for the registry package.
package testutil

import "github.com/ankeesler/btool/node/registry"

func Index() *registry.Index {
	return &registry.Index{
		Files: []registry.IndexFile{
			registry.IndexFile{
				Path:   "file_a_btool.yml",
				SHA256: "e3b5fa05ab6b80f6004e3e7bd88fae76e3b4dbaa21b696ecc619a532a1ac0875",
			},
			registry.IndexFile{
				Path:   "some/path/to/file_b_btool.yml",
				SHA256: "1d77c3b1aa9874e57e3ef410f584333ebc3f2a019727f8e29ad96deda8d8c8b5",
			},
		},
	}
}

func FileAGaggle() *registry.Gaggle {
	return &registry.Gaggle{
		Metadata: map[string]interface{}{
			"foo": "bar",
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
			"bat": "mayonnaise",
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
