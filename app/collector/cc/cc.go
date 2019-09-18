// Package cc provides collector.Producer's and collector.Consumer's specific to
// C/C++ code.
package cc

// Labels is a data structure that stores unmarshaled node.Node Label data.
// This data is common to all C/C++ related node.Node's.
type Labels struct {
	// Include paths required when compiling this thing.
	IncludePaths []string `mapstructure:"io.btool.collector.cc.includePaths"`

	// Libraries required when linking this thing.
	Libraries []string `mapstructure:"io.btool.collector.cc.libraries"`
}
