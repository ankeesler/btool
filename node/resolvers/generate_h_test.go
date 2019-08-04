package resolvers_test

import (
	"path/filepath"
	"testing"

	"github.com/ankeesler/btool/node"
	"github.com/ankeesler/btool/node/resolvers"
	"github.com/spf13/afero"
)

func TestGenerateH(t *testing.T) {
	data := []struct {
		file      string
		exContent string
	}{
		{
			file: "some/path/to/file.h",
			exContent: `#ifndef SOME_PATH_TO_FILE_H_
#define SOME_PATH_TO_FILE_H_

#endif // SOME_PATH_TO_FILE_H_
`,
		},
		{
			file: "some/path/to/file_with_underscores.h",
			exContent: `#ifndef SOME_PATH_TO_FILE_WITH_UNDERSCORES_H_
#define SOME_PATH_TO_FILE_WITH_UNDERSCORES_H_

#endif // SOME_PATH_TO_FILE_WITH_UNDERSCORES_H_
`,
		},
	}

	for _, datum := range data {
		t.Run(datum.file, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			dir := "some/dir"
			gh := resolvers.NewGenerateH(fs, dir)

			n := node.New(datum.file)
			if err := gh.Resolve(n); err != nil {
				t.Fatal(err)
			}

			data, err := afero.ReadFile(fs, filepath.Join(dir, datum.file))
			if err != nil {
				t.Fatal(err)
			}

			if ex, ac := datum.exContent, string(data); ex != ac {
				t.Error(ex, "!=", ac)
			}
		})
	}
}
