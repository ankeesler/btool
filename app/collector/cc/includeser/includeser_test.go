package includeser_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/ankeesler/btool/app/collector/cc/includeser"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestIncludeserIncludes(t *testing.T) {
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
		f, err := ioutil.TempFile("", "btool_includeser_test")
		require.Nil(t, err)
		defer func() {
			require.Nil(t, f.Close())
		}()

		_, err = fmt.Fprintln(f, datum.data)
		require.Nil(t, err)

		i := includeser.New(afero.NewOsFs())
		includes, err := i.Includes(f.Name())
		require.Nil(t, err)
		require.Equal(t, datum.includes, includes)
	}
}
