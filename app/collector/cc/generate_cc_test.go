package cc_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/ankeesler/btool/app/collector/cc"
	"github.com/spf13/afero"
)

func TestGenerate(t *testing.T) {
	data := []struct {
		root string
		name string

		path string

		failure bool

		ifndef    string
		namespace string
		class     string
	}{
		{
			root: "/tuna/root",
			name: "projectname",

			path: "basic",

			failure: false,

			ifndef:    "PROJECTNAME_BASIC_H_",
			namespace: "projectname",
			class:     "Basic",
		},
		{
			root: "/tuna/root",
			name: "projectname",

			path: "basic_snake_case",

			failure: false,

			ifndef:    "PROJECTNAME_BASICSNAKECASE_H_",
			namespace: "projectname",
			class:     "BasicSnakeCase",
		},
		{
			root: "/tuna/root",
			name: "dash-project-name",

			path: "basic_snake_case",

			failure: false,

			ifndef:    "DASHPROJECTNAME_BASICSNAKECASE_H_",
			namespace: "dashprojectname",
			class:     "BasicSnakeCase",
		},
		{
			root: "/tuna/root",
			name: "projectname",

			path: "file_ending_failure.h",

			failure: true,
		},
		{
			root: "/tuna/root",
			name: "projectname",

			path: "some/path/to/a/class",

			failure: false,

			ifndef:    "PROJECTNAME_SOME_PATH_TO_A_CLASS_H_",
			namespace: "projectname::some::path::to::a",
			class:     "Class",
		},
		{
			root: "/tuna/root",
			name: "projectname",

			path: "some/path/to/a/class_snake_case",

			failure: false,

			ifndef:    "PROJECTNAME_SOME_PATH_TO_A_CLASSSNAKECASE_H_",
			namespace: "projectname::some::path::to::a",
			class:     "ClassSnakeCase",
		},
		{
			path:    "some/path/to/a/file_ending_failure.h",
			failure: true,
		},
	}

	for _, datum := range data {
		fs := afero.NewMemMapFs()
		err := cc.GenerateClass(fs, datum.root, datum.name, datum.path)

		if err != nil {
			if !datum.failure {
				t.Errorf("%s: %s", datum.path, err.Error())
			}
			continue
		} else {
			if datum.failure {
				t.Errorf("%s: expected failure", datum.path)
				continue
			}
		}

		exHContent := fmt.Sprintf(`#ifndef %s
#define %s

namespace %s {

class %s {

};

};  // namespace %s

#endif  // %s
`,
			datum.ifndef,
			datum.ifndef,
			datum.namespace,
			datum.class,
			datum.namespace,
			datum.ifndef,
		)
		acHContent, err := afero.ReadFile(fs, filepath.Join(datum.root, datum.path+".h"))
		if err != nil {
			t.Errorf("%s: %s", datum.path, err.Error())
		} else if exHContent != string(acHContent) {
			t.Errorf(
				"%s: expected '%s', actual '%s'",
				datum.path,
				exHContent,
				string(acHContent),
			)
		}

		exCContent := fmt.Sprintf(`#include "%s"

namespace %s {

}; // namespace %s
`,
			datum.path+".h",
			datum.namespace,
			datum.namespace,
		)
		acCContent, err := afero.ReadFile(fs, filepath.Join(datum.root, datum.path+".cc"))
		if err != nil {
			t.Errorf("%s: %s", datum.path, err.Error())
		} else if exCContent != string(acCContent) {
			t.Errorf(
				"%s: expected '%s', actual '%s'",
				datum.path,
				exCContent,
				string(acCContent),
			)
		}
	}
}
