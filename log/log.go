package log

import (
	"fmt"
	"io"
	"os"
)

type Level int

const (
	ERROR Level = iota
	INFO
)

type Logger struct {
	Level Level
	out   io.Writer
}

var Default = New(ERROR, os.Stderr)

func New(level Level, out io.Writer) *Logger {
	return &Logger{
		Level: level,
		out:   out,
	}
}

func (l *Logger) Log(level Level, message string) (int, error) {
	if l.Level < level {
		return 0, nil
	}
	return fmt.Fprintln(l.out, message)
}

func Log(level Level, message string) (int, error) {
	return Default.Log(level, message)
}

func (l *Logger) Logf(level Level, format string, a ...interface{}) (int, error) {
	if l.Level < level {
		return 0, nil
	}
	return fmt.Fprintf(l.out, format, a...)
}

func Logf(level Level, format string, a ...interface{}) (int, error) {
	return Default.Logf(level, format, a...)
}

func (l *Logger) Error(s string) (int, error) {
	return l.Log(ERROR, s)
}

func Error(s string) (int, error) {
	return Default.Error(s)
}

func (l *Logger) Errorf(s string, a ...interface{}) (int, error) {
	return l.Logf(ERROR, s, a...)
}

func Errorf(s string, a ...interface{}) (int, error) {
	return Default.Errorf(s, a...)
}

func (l *Logger) Info(s string) (int, error) {
	return l.Log(INFO, s)
}

func Info(s string) (int, error) {
	return Default.Info(s)
}

func (l *Logger) Infof(s string, a ...interface{}) (int, error) {
	return l.Logf(INFO, s, a...)
}

func Infof(s string, a ...interface{}) (int, error) {
	return Default.Infof(s, a...)
}
