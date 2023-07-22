package logger

import (
	"os"
	"strings"
)

type Level int8

const (
	DebugLevel Level = iota << 1
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

var (
	defaultLevel = DebugLevel
	allLevel     = []Level{DebugLevel, InfoLevel, WarnLevel, ErrorLevel, FatalLevel}
)

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return ""
	}
}

func DefaultLevel() Level {
	return fromString(os.Getenv("LOGGER_LEVEL"))
}

func fromString(level string) Level {
	if level != "" {
		levelUpper := strings.ToUpper(level)

		for _, l := range allLevel {
			if levelUpper == l.String() {
				return l
			}
		}
	}

	return defaultLevel
}
