package testutil

import (
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
)

type TestingFormatter struct {
	t *testing.T
}

func NewTestingFormatter(t *testing.T) *TestingFormatter {
	return &TestingFormatter{
		t: t,
	}
}

func (tf *TestingFormatter) Format(entry *logrus.Entry) ([]byte, error) {
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
