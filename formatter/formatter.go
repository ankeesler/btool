// Package formatter provides a custom logrus formatter for btool.
package formatter

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

type formatter struct {
}

func New() logrus.Formatter {
	return &formatter{}
}

func (f *formatter) Format(entry *logrus.Entry) ([]byte, error) {
	s := fmt.Sprintf(
		"%s | %s | %s\n",
		color.CyanString("btool"),
		colorLevel(entry.Level),
		color.HiBlackString(entry.Message),
	)
	return []byte(s), nil
}

func colorLevel(level logrus.Level) string {
	switch level {
	case logrus.DebugLevel:
		return color.HiBlackString("%7s", level)
	case logrus.InfoLevel:
		return color.HiBlueString("%7s", level)
	case logrus.WarnLevel:
		return color.YellowString("%7s", level)
	case logrus.ErrorLevel:
		return color.HiRedString("%7s", level)
	case logrus.FatalLevel:
		return color.HiRedString("%7s", level)
	default:
		return fmt.Sprintf("%7s", level)
	}
}
