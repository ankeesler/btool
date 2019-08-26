// Package testutil provides test utilities for the registry package.
package testutil

import "github.com/ankeesler/btool/registry"

func Index() *registry.Index {
	return &registry.Index{
		Name: "some-index",
		Files: []registry.IndexFile{
			registry.IndexFile{
				Path:   "file_a_btool.yml",
				SHA256: "bcd28f7fa81b695eea671c474c11d842410e1925f794d681ce85f9bd4675c0c0",
			},
			registry.IndexFile{
				Path:   "some/path/to/file_b_btool.yml",
				SHA256: "5d0418f0a8c7f380bc0fb1766d5766a2d2f5f468e552113c4242b1637597d781",
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
					Name: "tuna-resolver",
					Config: map[string]interface{}{
						"tuna": "tuna",
					},
				},
			},
			&registry.Node{
				Name:         "fish",
				Dependencies: []string{"tuna"},
				Resolver: registry.Resolver{
					Name: "fish-resolver",
					Config: map[string]interface{}{
						"fish": "fish",
					},
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
					Name: "marlin-resolver",
					Config: map[string]interface{}{
						"marlin": "marlin",
					},
				},
			},
			&registry.Node{
				Name:         "bacon",
				Dependencies: []string{"marlin"},
				Resolver: registry.Resolver{
					Name: "bacon-resolver",
					Config: map[string]interface{}{
						"bacon": "bacon",
					},
				},
			},
		},
	}
}
