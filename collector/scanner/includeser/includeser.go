// Package includeser provides the ability to parse #include directives from a
// C/C++ file.
package includeser

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
)

type Includeser struct {
}

func New() *Includeser {
	return &Includeser{}
}

func (i *Includeser) Includes(path string) ([]string, error) {
	bytes, err := ioutil.ReadFile(path)
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
