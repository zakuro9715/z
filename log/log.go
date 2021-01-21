package log

import (
	"fmt"
	"io"
	"os"
)

type Level int

const (
	FATAL Level = iota
	INFO
)

type Logger struct {
	Level Level
	out   io.Writer
}

var Default = New(FATAL, os.Stderr)

func New(level Level, out io.Writer) *Logger {
	return &Logger{
		Level: level,
		out:   out,
	}
}

func (l *Logger) Log(level Level, format string, a ...interface{}) (int, error) {
	if l.Level < level {
		return 0, nil
	}
	return fmt.Fprintf(l.out, format+"\n", a...)
}

func Log(level Level, format string, a ...interface{}) (int, error) {
	return Default.Log(level, format, a...)
}

func (l *Logger) Fatal(s string, a ...interface{}) {
	l.Log(FATAL, "[FATAL] "+s, a...)
	os.Exit(2)
}

func Fatal(s string, a ...interface{}) {
	Default.Fatal(s, a...)
}

func (l *Logger) Info(s string, a ...interface{}) (int, error) {
	return l.Log(INFO, "[INFO] "+s, a...)
}

func Info(s string, a ...interface{}) (int, error) {
	return Default.Info(s, a...)
}
