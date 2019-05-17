package includes_test

import (
	"reflect"
	"testing"

	"github.com/ankeesler/btool/scanner/includes"
)

func TestScan(t *testing.T) {
	data := []struct {
		name     string
		data     string
		includes []string
	}{
		{
			name: "basic",
			data: `
#include "some/path/to/file.h"

#define IGNORE_THIS "hey"
#include "some/other/path/to/file.h"

#include "tuna.h"
#   include "weird-spaces.h"
#include"no-spaces.h"
`,
			includes: []string{
				"some/path/to/file.h",
				"some/other/path/to/file.h",
				"tuna.h",
				"weird-spaces.h",
				"no-spaces.h",
			},
		},
	}

	for _, datum := range data {
		includes, err := includes.Parse([]byte(datum.data))
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(datum.includes, includes) {
			t.Errorf("%s: %s != %s", datum.name, datum.includes, includes)
		}
	}
}
