package registryapi_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankeesler/btool/registry/registryapi"
	"github.com/ankeesler/btool/registry/registryapi/registryapifakes"
)

func TestRegistryAPI(t *testing.T) {
	data := []struct {
		rFunc func() *registryapifakes.FakeRegistry

		path string
		code int
		body string
	}{
		{
			rFunc: func() *registryapifakes.FakeRegistry {
				index := "index"
				indexBuf := bytes.NewBuffer([]byte(index))

				r := &registryapifakes.FakeRegistry{}
				r.GetReturnsOnCall(0, indexBuf)
				return r
			},

			path: "/",
			code: http.StatusOK,
			body: "index",
		},
		{
			rFunc: func() *registryapifakes.FakeRegistry {
				file := "file"
				fileBuf := bytes.NewBuffer([]byte(file))

				r := &registryapifakes.FakeRegistry{}
				r.GetReturnsOnCall(0, fileBuf)
				return r
			},

			path: "/path/to/file.yml",
			code: http.StatusOK,
			body: "file",
		},
		{
			rFunc: func() *registryapifakes.FakeRegistry {
				r := &registryapifakes.FakeRegistry{}
				r.GetReturnsOnCall(0, nil)
				return r
			},

			path: "/notfound",
			code: http.StatusNotFound,
			body: "",
		},
	}
	for _, datum := range data {
		r := datum.rFunc()
		api := registryapi.New(r)

		req := httptest.NewRequest(
			http.MethodGet,
			datum.path,
			nil,
		)
		rsp := httptest.NewRecorder()

		api.ServeHTTP(rsp, req)
		if ex, ac := datum.code, rsp.Code; ex != ac {
			t.Error(ex, "!=", ac)
		}

		if datum.code != http.StatusNotFound {
			if ex, ac := datum.body, rsp.Body.String(); ex != ac {
				t.Error(ex, "!=", ac)
			}
		}

		if ex, ac := 1, r.GetCallCount(); ex != ac {
			t.Error(ex, "!=", ac)
		}
	}
}
