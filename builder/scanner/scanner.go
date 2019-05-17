package scanner

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/pkg/errors"
)

type Scanner struct {
}

func New() *Scanner {
	return &Scanner{}
}

func (s *Scanner) Scan(data []byte) ([]string, error) {
	includes := make([]string, 0, 2)

	realScanner := bufio.NewScanner(bytes.NewBuffer(data))
	for realScanner.Scan() {
		line := strings.Replace(realScanner.Text(), " ", "", -1)
		if strings.HasPrefix(line, "#include\"") {
			include := strings.TrimPrefix(
				strings.TrimSuffix(line, "\""),
				"#include\"",
			)
			includes = append(includes, include)
		}
	}

	if realScanner.Err() != nil {
		return nil, errors.Wrap(realScanner.Err(), "scan")
	}

	return includes, nil
}
