// Package includeser provides the ability to parse #include directives from a
// C/C++ file.
package includeser

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type Includeser struct {
	fs afero.Fs
}

func New(fs afero.Fs) *Includeser {
	return &Includeser{
		fs: fs,
	}
}

func (i *Includeser) Includes(path string) ([]string, error) {
	bytes, err := afero.ReadFile(i.fs, path)
	if err != nil {
		return nil, errors.Wrap(err, "read file")
	}

	includes, err := parse(bytes)
	if err != nil {
		return nil, errors.Wrap(err, "parse")
	}

	return includes, nil
}

func parse(data []byte) ([]string, error) {
	includes := make([]string, 0, 2)

	scanner := bufio.NewScanner(bytes.NewBuffer(data))
	for scanner.Scan() {
		line := strings.Replace(scanner.Text(), " ", "", -1)
		if strings.HasPrefix(line, "#include\"") {
			include := strings.TrimPrefix(
				strings.TrimSuffix(line, "\""),
				"#include\"",
			)
			includes = append(includes, include)
		}
	}

	if scanner.Err() != nil {
		return nil, errors.Wrap(scanner.Err(), "scan")
	}

	return includes, nil
}
