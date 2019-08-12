// Package api returns a http.Handler that will serve the btool registry.
package api

import (
	"fmt"
	"net/http"

	"github.com/ankeesler/btool/node/registry"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Registry

// Registry is an object that can retrieve registry.Node's.
type Registry interface {
	// Index should return the registry.Index associated with this particular
	// Registry. If any error occurs, an error should be returned.
	Index() (*registry.Index, error)
	// Nodes should return the registry.Node's associated with the provided
	// registry.IndexFile.Path. If any error occurs, an error should be returned.
	// If no registry.Node's exist for the provided string, then an empty
	// slice should be returned.
	Nodes(string) ([]*registry.Node, error)
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

	logrus.Debugf("handling %s %s", req.Method, path)

	rsp.Header().Set("Content-Type", "application/x-yaml")

	var object interface{}
	var err error
	if path == "" {
		object, err = ra.r.Index()
	} else {
		var nodes []*registry.Node
		nodes, err = ra.r.Nodes(path)
		logrus.Debugf("nodes returned %s", nodes)

		if len(nodes) != 0 {
			object = nodes
		}
	}
	logrus.Debug("response object: ", object)

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

	logrus.Debugf("encoding object %s", object)
	if err := e.Encode(object); err != nil {
		rsp.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rsp, "---\nerror: %s\n", err.Error())
		return
	}
}
