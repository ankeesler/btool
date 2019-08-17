// Package api returns a http.Handler that will serve the btool registry.
package api

import (
	"fmt"
	"net/http"

	"github.com/ankeesler/btool/log"
	"github.com/ankeesler/btool/node/registry"
	"gopkg.in/yaml.v2"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Registry

// Registry is an object that can retrieve registry.Gaggle's.
type Registry interface {
	// Index should return the registry.Index associated with this particular
	// Registry. If any error occurs, an error should be returned.
	Index() (*registry.Index, error)
	// Gaggle should return the registry.Gaggle associated with the provided
	// registry.IndexFile.Path. If any error occurs, an error should be returned.
	// If no registry.Gaggle exists for the provided string, then nil, nil should
	// be returned.
	Gaggle(string) (*registry.Gaggle, error)
}

type registryApi struct {
	r Registry
}

// New returns a new http.Handler that serves a btool registry.
func New(r Registry) http.Handler {
	return &registryApi{
		r: r,
	}
}

func (ra *registryApi) ServeHTTP(rsp http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	if path[0] == '/' {
		path = path[1:]
	}

	log.Debugf("handling %s %s", req.Method, path)

	rsp.Header().Set("Content-Type", "application/x-yaml")

	var object interface{}
	var err error
	if path == "" {
		object, err = ra.r.Index()
	} else {
		var gaggle *registry.Gaggle
		gaggle, err = ra.r.Gaggle(path)
		if gaggle != nil {
			object = gaggle
		}
	}
	log.Debugf("response object: %s", object)

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

	log.Debugf("encoding object %s", object)
	if err := e.Encode(object); err != nil {
		rsp.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rsp, "---\nerror: %s\n", err.Error())
		return
	}
}
