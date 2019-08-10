package registry

// Index describes the files in the registry.
type Index struct {
	Files []IndexFile `yaml:"files"`
}

// IndexFile describes a single file in the registry.
type IndexFile struct {
	Path   string `yaml:"path"`
	SHA256 string `yaml:"sha256"`
}
