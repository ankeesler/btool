package cc

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/ankeesler/btool/log"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// GenerateClass will generate a .cc and .h file with the provided project of the
// C++ class at the filepath.Join(root, path) prefix. The name is dictated by the
// provided path.
func GenerateClass(fs afero.Fs, root, project, path string) error {
	if filepath.Ext(path) != "" {
		return errors.New("do not put extension onto generate path (e.g., some/path/to/file)")
	}

	hFile := filepath.Join(root, path+".h")
	hContent := hContent(project, path)
	if err := fs.MkdirAll(filepath.Dir(hFile), 0755); err != nil {
		return errors.Wrap(err, "mkdir all")
	}
	if err := afero.WriteFile(fs, hFile, hContent, 0600); err != nil {
		return errors.Wrap(err, "write "+hFile)
	}
	log.Infof("generated %s", hFile)

	cFile := filepath.Join(root, path+".cc")
	cContent := cContent(project, path)
	if err := fs.MkdirAll(filepath.Dir(cFile), 0755); err != nil {
		return errors.Wrap(err, "mkdir all")
	}
	if err := afero.WriteFile(fs, cFile, cContent, 0600); err != nil {
		return errors.Wrap(err, "write "+cFile)
	}
	log.Infof("generated %s", cFile)

	return nil
}

func hContent(project, path string) []byte {
	buf := bytes.NewBuffer([]byte{})

	ifndef := ifndef(project, path)
	namespace := namespace(project, path)
	class := class(path)

	fmt.Fprintf(
		buf,
		`#ifndef %s
#define %s

namespace %s {

class %s {

};

}; // namespace %s

#endif // %s
`,
		ifndef,
		ifndef,
		namespace,
		class,
		namespace,
		ifndef,
	)

	return buf.Bytes()
}

func cContent(project, path string) []byte {
	buf := bytes.NewBuffer([]byte{})

	namespace := namespace(project, path)

	fmt.Fprintf(
		buf,
		`#include "%s"

namespace %s {

}; // namespace %s
`,
		path+".h",
		namespace,
		namespace,
	)

	return buf.Bytes()
}

func ifndef(project, path string) string {
	buf := bytes.NewBuffer([]byte{})

	fmt.Fprintf(buf, "%s_", strings.ToUpper(cleanProjectName(project)))
	fmt.Fprintf(
		buf,
		"%s_H_",
		strings.NewReplacer(
			"_", "",
			"/", "_",
		).Replace(strings.ToUpper(path)),
	)

	return buf.String()
}

func namespace(project, path string) string {
	buf := bytes.NewBuffer([]byte{})

	buf.WriteString(cleanProjectName(project))

	dir := filepath.Dir(path)
	if dir != "." {
		fmt.Fprintf(buf, "::%s", strings.ReplaceAll(dir, "/", "::"))
	}

	return buf.String()
}

func class(path string) string {
	base := []byte(filepath.Base(path))
	class := make([]rune, 0)

	upper := true
	for i := range base {
		r := rune(base[i])
		if r == '_' {
			upper = true
		} else if upper {
			class = append(class, unicode.ToUpper(r))
			upper = false
		} else {
			class = append(class, r)
		}
	}

	//from := 0
	//for {
	//	to := bytes.IndexByte(base[from:], '_')
	//	if to == -1 {
	//		break
	//	}
	//
	//	copy(class, base[from:to])
	//	from = to
	//}

	return string(class)
}

func cleanProjectName(project string) string {
	return strings.NewReplacer(
		"-", "",
	).Replace(project)
}
