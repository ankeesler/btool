// Package registryapi returns a http.Handler that will serve the btool registry.
package registryapi

import (
	"fmt"
	"net/http"

	"github.com/ankeesler/btool/registry"
	"gopkg.in/yaml.v2"
)

type registryApi struct {
	r registry.Registry
}

// New returns a new http.Handler that serves a btool registry.Registry.
func New(r registry.Registry) http.Handler {
	return &registryApi{
		r: r,
	}
}

func (ra *registryApi) ServeHTTP(rsp http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	rsp.Header().Set("Content-Type", "application/x-yaml")

	var object interface{}
	var err error
	if path == "/" {
		object, err = ra.r.Index()
	} else {
		object, err = ra.r.Nodes(path)
	}

	if err != nil {
		rsp.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rsp, "---\nerror: %s\n", err.Error())
		return
	}

	if object == nil {
		rsp.WriteHeader(http.StatusNotFound)
		return
	}

	e := yaml.NewEncoder(rsp)
	defer e.Close()

	if err := e.Encode(object); err != nil {
		rsp.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rsp, "---\nerror: %s\n", err.Error())
		return
	}
}
