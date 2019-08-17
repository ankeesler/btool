// Package log provides a logger for the btool project.
package log

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/fatih/color"
)

// Level is an enum that describes how verbose a log should be.
type Level int

// The Level* constants define how verbose the log should be.
const (
	LevelDebug Level = iota
	LevelInfo
	LevelError
)

// CurrentLevel is the current Level of logging.
var CurrentLevel = LevelDebug

// ParseLevel will parse a Level from a string.
func ParseLevel(level string) (Level, error) {
	switch level {
	case "debug":
		return LevelDebug, nil
	case "info":
		return LevelInfo, nil
	default:
		return 0, errors.New("unknown level: " + level)
	}
}

// Errorf will print something to the log with LevelError.
func Errorf(format string, args ...interface{}) {
	log(LevelError, format, args...)
}

// Infof will print something to the log with LevelInfo.
func Infof(format string, args ...interface{}) {
	log(LevelInfo, format, args...)
}

// Debugf will print something to the log with LevelDebug.
func Debugf(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}

	format = fmt.Sprintf("%10s:%3d | %s", filepath.Base(file), line, format)
	log(LevelDebug, format, args...)
}

func log(level Level, format string, args ...interface{}) {
	if level >= CurrentLevel {
		fmt.Printf(
			"%s | %s | %s\n",
			color.CyanString("btool"),
			colorLevel(level),
			color.HiBlackString(fmt.Sprintf(format, args...)),
		)
	}
}

func colorLevel(level Level) string {
	switch level {
	case LevelDebug:
		return color.HiBlackString("%5s", "debug")
	case LevelInfo:
		return color.HiBlueString("%5s", "info")
	default:
		return fmt.Sprintf("%5s", "???")
	}
}
