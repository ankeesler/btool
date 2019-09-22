// +build tools

package tools

import (
	_ "github.com/golang/protobuf/protoc-gen-go"
	_ "github.com/maxbrunsfeld/counterfeiter/v6"
	_ "google.golang.org/grpc"
)

// This file imports packages that are used when running go generate, or used
// during the development process but not otherwise depended on by built code.
