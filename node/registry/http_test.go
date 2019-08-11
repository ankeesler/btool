package registry_test

import (
	"bytes"
	"net/http"
	"strings"
	"testing"

	"github.com/ankeesler/btool/node/registry"
	"github.com/ankeesler/btool/node/registry/registryfakes"
	"github.com/ankeesler/btool/node/registry/testutil"
	"github.com/go-test/deep"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type readCloser struct {
	buf *bytes.Buffer
}

func (rcb *readCloser) Read(b []byte) (int, error) {
	return rcb.buf.Read(b)
}

func (rcb *readCloser) Close() error {
	return nil
}

func TestHTTPRegistryIndex(t *testing.T) {
	exI := testutil.Index()

	buf200 := bytes.NewBuffer([]byte{})
	e := yaml.NewEncoder(buf200)
	defer e.Close()
	if err := e.Encode(exI); err != nil {
		t.Fatal(err)
	}

	rsp200 := http.Response{
		StatusCode: http.StatusOK,
		Body:       &readCloser{buf200},
	}
	rsp404 := http.Response{
		StatusCode: http.StatusNotFound,
		Body:       &readCloser{buf200},
	}
	rsp500 := http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       &readCloser{buf200},
	}

	c := &registryfakes.FakeHTTPClient{}
	c.GetReturnsOnCall(0, &rsp200, nil)
	c.GetReturnsOnCall(1, &rsp404, nil)
	c.GetReturnsOnCall(2, &rsp500, nil)
	c.GetReturnsOnCall(3, nil, errors.New("some error"))

	r := registry.NewHTTPRegistry("some url", c)

	// 200
	acI, err := r.Index()
	if err != nil {
		t.Fatal(err)
	}

	if diff := deep.Equal(exI, acI); diff != nil {
		t.Error(exI, "!=", acI)
	}

	// 404
	_, err = r.Index()
	if err == nil {
		t.Error("expected error")
	}

	// 500
	_, err = r.Index()
	if err == nil {
		t.Error("expected error")
	}

	// error
	_, err = r.Index()
	if err == nil {
		t.Error("expected error")
	} else if !strings.Contains(err.Error(), "some error") {
		t.Error()
	}

	if ex, ac := 4, c.GetCallCount(); ex != ac {
		t.Error(ex, "!=", ac)
	}

	for i := 0; i < c.GetCallCount(); i++ {
		url := c.GetArgsForCall(i)
		if ex, ac := "some url", url; ex != ac {
			t.Error(i, "->", ex, "!=", ac)
		}
	}
}

func andTestHTTPRegistryNodes(t *testing.T) {
	exNodes := testutil.Nodes()

	buf200 := bytes.NewBuffer([]byte{})
	e := yaml.NewEncoder(buf200)
	defer e.Close()
	if err := e.Encode(exNodes); err != nil {
		t.Fatal(err)
	}

	rsp200 := http.Response{
		StatusCode: http.StatusOK,
		Body:       &readCloser{buf200},
	}
	rsp404 := http.Response{
		StatusCode: http.StatusNotFound,
		Body:       &readCloser{buf200},
	}
	rsp500 := http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       &readCloser{buf200},
	}

	c := &registryfakes.FakeHTTPClient{}
	c.GetReturnsOnCall(0, &rsp200, nil)
	c.GetReturnsOnCall(1, &rsp404, nil)
	c.GetReturnsOnCall(2, &rsp500, nil)
	c.GetReturnsOnCall(3, nil, errors.New("some error"))

	r := registry.NewHTTPRegistry("some url", c)

	// 200
	acNodes, err := r.Nodes("some/nodes")
	if err != nil {
		t.Fatal(err)
	}

	if diff := deep.Equal(exNodes, acNodes); diff != nil {
		t.Error(exNodes, "!=", acNodes)
	}

	// 404
	acNodes, err = r.Nodes("some/nodes")
	if err != nil {
		t.Error(err)
	} else if acNodes != nil {
		t.Error("expected nil nodes")
	}

	// 500
	_, err = r.Nodes("some/nodes")
	if err == nil {
		t.Error("expected error")
	}

	// error
	_, err = r.Nodes("some/nodes")
	if err == nil {
		t.Error("expected error")
	} else if !strings.Contains(err.Error(), "some error") {
		t.Error()
	}

	if ex, ac := 4, c.GetCallCount(); ex != ac {
		t.Error(ex, "!=", ac)
	}

	for i := 0; i < c.GetCallCount(); i++ {
		url := c.GetArgsForCall(i)
		if ex, ac := "some url/some/nodes", url; ex != ac {
			t.Error(i, "->", ex, "!=", ac)
		}
	}
}
