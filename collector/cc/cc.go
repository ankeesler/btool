package cc

// The following are Labels used to decorate node.Node's.
const (
	// Include paths required when compiling this thing.
	// A list of comma separated paths.
	LabelIncludePaths = "io.btool.cc.includePaths"

	// Libraries required when linking this thing.
	// A list of comma separated paths.
	LabelLibraries = "io.btool.cc.libraries"
)
