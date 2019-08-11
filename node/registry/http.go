package registry

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . HTTPClient

// HTTPClient is an object that can send a GET HTTP request.
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

// HTTPRegistry retrieves Index()/Node() information from an HTTP URL.
type HTTPRegistry struct {
	url        string
	httpClient HTTPClient
}

var errNotFound = errors.New("not found")

// NewHTTPRegistry returns a HTTPRegistry at some URL.
func NewHTTPRegistry(url string, httpClient HTTPClient) *HTTPRegistry {
	return &HTTPRegistry{
		url:        url,
		httpClient: httpClient,
	}
}

func (hr *HTTPRegistry) Index() (*Index, error) {
	i := new(Index)
	if err := hr.get(hr.url, i); err != nil {
		return nil, errors.Wrap(err, "get")
	}
	return i, nil
}

func (hr *HTTPRegistry) Nodes(path string) ([]*Node, error) {
	nodes := make([]*Node, 0)
	if err := hr.get(hr.url+path, &nodes); err != nil {
		if err == errNotFound {
			return nil, nil
		} else {
			return nil, errors.Wrap(err, "get")
		}
	}
	return nodes, nil
}

func (hr *HTTPRegistry) get(url string, object interface{}) error {
	rsp, err := hr.httpClient.Get(url)
	if err != nil {
		return errors.Wrap(err, "http get")
	}
	defer rsp.Body.Close()

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return errors.Wrap(err, "read body")
	}

	if rsp.StatusCode == http.StatusNotFound {
		return errNotFound
	} else if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", rsp.StatusCode)
	}

	if err := yaml.Unmarshal(data, object); err != nil {
		return errors.Wrap(err, "unmarshal")
	}

	return nil
}
