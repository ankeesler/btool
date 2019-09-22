// Package api provides the btool node.Node API.
//
// When wanting to regenerate the API protobuf definitions, you should run:
//   go generate github.com/ankeesler/btool/node/api
package api

//go:generate protoc --go_out=plugins=grpc:. v1/node.proto
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 v1 RegistryClient
