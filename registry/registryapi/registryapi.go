// Package registryapi returns a http.Handler that will serve the btool registry.
package registryapi

import (
	"fmt"
	"io"
	"net/http"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Registry

// Registry returns io.Reader's for the files in a btool registry.
type Registry interface {
	Get(string) io.Reader
}

type registryApi struct {
	r Registry
}

// New returns a new http.Handler that serves the btool registry from a given
// directory of *_btool.yml files.
func New(r Registry) http.Handler {
	return &registryApi{
		r: r,
	}
}

func (ra *registryApi) ServeHTTP(rsp http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	if path == "/" {
		path = "index.yml"
	}

	file := ra.r.Get(path)
	if file == nil {
		rsp.WriteHeader(http.StatusNotFound)
		return
	}

	if _, err := io.Copy(rsp, file); err != nil {
		rsp.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rsp, "error: %s\n", err.Error())
		return
	}

	rsp.Header().Set("Content-Type", "application/x-yaml")
}
