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
		color.BlackString(entry.Message),
	)
	return []byte(s), nil
}

func colorLevel(level logrus.Level) string {
	switch level {
	case logrus.DebugLevel:
		return color.HiBlackString("%5s", level)
	case logrus.InfoLevel:
		return color.BlueString("%5s", level)
	case logrus.WarnLevel:
		return color.YellowString("%5s", level)
	case logrus.ErrorLevel:
		return color.HiRedString("%5s", level)
	case logrus.FatalLevel:
		return color.RedString("%5s", level)
	default:
		return color.BlackString("%5s", level)
	}
}
