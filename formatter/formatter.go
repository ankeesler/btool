// Package formatter provides a custom logrus formatter for btool.
package formatter

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type formatter struct {
}

func New() logrus.Formatter {
	return &formatter{}
}

func (f *formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var file string
	var line int
	if entry.HasCaller() {
		file = entry.Caller.File
		line = entry.Caller.Line
	} else {
		file = "?"
		line = 0
	}

	s := fmt.Sprintf(
		"btool | %s | %s:%d | %s\n",
		entry.Level,
		file,
		line,
		entry.Message,
	)
	return []byte(s), nil
}
