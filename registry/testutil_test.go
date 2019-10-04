package registry_test

import "github.com/ankeesler/btool/registry"

func index() *registry.Index {
	return &registry.Index{
		Files: []registry.IndexFile{
			registry.IndexFile{
				Path:   "file_a_btool.yml",
				SHA256: "f66a9025bf3d41dad609ea5b90ffd0e8d72b50736f416a6c19e505de5c7f3129",
			},
			registry.IndexFile{
				Path:   "some/path/to/file_b_btool.yml",
				SHA256: "ec97db25b5a77e6b0f5ac8017d6125bd751e2a1d89b8463b7e8a3cc0b2300daa",
			},
		},
	}
}

func fileAGaggle() *registry.Gaggle {
	return &registry.Gaggle{
		Metadata: map[string]interface{}{
			"foo":     "bar",
			"project": "some-project",
		},
		Nodes: []*registry.Node{
			&registry.Node{
				Name:         "tuna",
				Dependencies: []string{},
				Labels:       map[string]interface{}{},
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
				Labels:       map[string]interface{}{},
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

func fileBGaggle() *registry.Gaggle {
	return &registry.Gaggle{
		Metadata: map[string]interface{}{
			"bat":     "mayonnaise",
			"project": "some-other-project",
		},
		Nodes: []*registry.Node{
			&registry.Node{
				Name:         "marlin",
				Dependencies: []string{},
				Labels:       map[string]interface{}{},
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
				Labels:       map[string]interface{}{},
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
