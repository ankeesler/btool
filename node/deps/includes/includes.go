// Package includes provides the ability to parse #include directives from a
// C/C++ file.
package includes

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/pkg/errors"
)

func Parse(data []byte) ([]string, error) {
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
