package registry_test

import (
	"testing"

	"github.com/ankeesler/btool/registry"
	"github.com/ankeesler/btool/registry/registryfakes"
)

func TestHTTPRegistry(t *testing.T) {
	t.Fatal("WRITE ME")

	data := []struct {
		url   string
		cFunc func() *registryfakes.FakeHTTPClient
	}{}
	for _, datum := range data {
		c := datum.cFunc()
		r := registry.NewHTTPRegistry(datum.url, c)
		_ = r
	}
}
