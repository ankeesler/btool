package registry

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ankeesler/btool/log"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . HTTPClient

// HTTPClient is an object that can send a GET HTTP request.
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

// HTTPRegistry retrieves Index/Gaggle from an HTTP/HTTPS URL.
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
	log.Debugf("index -> %s", i)
	return i, nil
}

func (hr *HTTPRegistry) Gaggle(path string) (*Gaggle, error) {
	gaggle := new(Gaggle)
	if err := hr.get(hr.url+"/"+path, gaggle); err != nil {
		if err == errNotFound {
			return nil, nil
		} else {
			return nil, errors.Wrap(err, "get")
		}
	}
	log.Debugf("Gaggle(%s) -> %s", path, gaggle)
	return gaggle, nil
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
	log.Debugf("get returned %d/%s", rsp.StatusCode, string(data))

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
