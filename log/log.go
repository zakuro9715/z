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

func (l *Logger) Log(level Level, message string) (int, error) {
	return l.Logf(level, message+"\n")
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

func (l *Logger) Fatal(s string) {
	l.Fatalf(s + "\n")
}

func Fatal(s string) {
	Default.Fatal(s)
}

func (l *Logger) Fatalf(s string, a ...interface{}) {
	l.Logf(FATAL, s, a...)
	os.Exit(2)
}

func Fatalf(s string, a ...interface{}) {
	Default.Fatalf(s, a...)
}

func (l *Logger) Info(s string) (int, error) {
	return l.Infof(s + "\n")
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
