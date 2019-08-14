// Package toolchain helps find the compiler/linker tools for btool to use.
package toolchain

import (
	"errors"
	"runtime"
)

// Toolchain represents the compilers and linkers needed for btool to run.
type Toolchain struct {
	CompilerC  string
	CompilerCC string
	Archiver   string
	Linker     string
}

var (
	linux = Toolchain{
		CompilerC:  "gcc",
		CompilerCC: "g++",
		Archiver:   "ar",
		Linker:     "gcc",
	}
	darwin = Toolchain{
		CompilerC:  "clang",
		CompilerCC: "clang++",
		Archiver:   "ar",
		Linker:     "clang",
	}
)

// Find returns the Toolchain object for this system, or an error if it fails.
func Find() (*Toolchain, error) {
	switch runtime.GOOS {
	case "linux":
		return &linux, nil
	case "darwin":
		return &darwin, nil
	default:
		return nil, errors.New("unknown toolchain for system")
	}
}
