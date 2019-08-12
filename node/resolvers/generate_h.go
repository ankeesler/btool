package resolvers

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type generateH struct {
	fs  afero.Fs
	dir string
}

// NewGenerateH returns a node.Resolver that generates a .h header file.
func NewGenerateH(fs afero.Fs, dir string) node.Resolver {
	return &generateH{
		fs:  fs,
		dir: dir,
	}
}

func (gh *generateH) Resolve(n *node.Node) error {
	ifndef := makeIfndef(n.Name)
	content := fmt.Sprintf(
		"#ifndef %s\n#define %s\n\n#endif // %s\n",
		ifndef,
		ifndef,
		ifndef,
	)

	file := filepath.Join(gh.dir, n.Name)
	if err := gh.fs.MkdirAll(filepath.Dir(file), 0755); err != nil {
		return errors.Wrap(err, fmt.Sprintf("mkdir (%s)", gh.dir))
	}

	if err := afero.WriteFile(
		gh.fs,
		file,
		[]byte(content),
		0644,
	); err != nil {
		return errors.Wrap(err, fmt.Sprintf("write file (%s)", file))
	}

	return nil
}

func makeIfndef(file string) string {
	buf := bytes.NewBuffer([]byte{})

	fmt.Fprintf(
		buf,
		"%s_H_",
		strings.NewReplacer(
			"/", "_",
			".H", "",
		).Replace(strings.ToUpper(file)),
	)

	return buf.String()
}
