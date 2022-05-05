package logger

import (
	"fmt"
	"os"
)

type Level int

const (
	Lsilent Level = iota
	Linfo
	Ldebug
	Ltrace
)

func (s Level) String() string {
	switch s {
	case Linfo:
		return "I"
	case Ldebug:
		return "D"
	case Ltrace:
		return "T"
	default:
		return ""
	}
}

var (
	logLevel = Linfo
)

func GetLevel() Level      { return logLevel }
func SetLevel(level Level) { logLevel = level }

func output(level Level, format string, v ...any) {
	if level <= logLevel {
		fmt.Fprintf(os.Stderr, "[fcli][%s] %s\n", level, fmt.Sprintf(format, v...))
	}
}

func Info(format string, v ...any) { output(Linfo, format, v...) }

func Debug(format string, v ...any) { output(Ldebug, format, v...) }

func Trace(format string, v ...any) { output(Ltrace, format, v...) }
