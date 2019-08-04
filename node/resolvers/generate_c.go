package resolvers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ankeesler/btool/node"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type generateC struct {
	fs  afero.Fs
	dir string
}

// NewGenerateC returns a node.Resolver that generates a .c source file.
func NewGenerateC(fs afero.Fs, dir string) node.Resolver {
	return &generateC{
		fs:  fs,
		dir: dir,
	}
}

func (gc *generateC) Resolve(n *node.Node) error {
	include := strings.ReplaceAll(filepath.Base(n.Name), ".c", ".h")
	content := fmt.Sprintf("#include \"%s\"\n", include)

	file := filepath.Join(gc.dir, n.Name)

	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		return errors.Wrap(err, fmt.Sprintf("mkdir (%s)", gc.dir))
	}

	if err := afero.WriteFile(
		gc.fs,
		file,
		[]byte(content),
		0644,
	); err != nil {
		return errors.Wrap(err, fmt.Sprintf("write file (%s)", file))
	}

	return nil
}
