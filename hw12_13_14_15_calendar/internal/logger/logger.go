package logger

import (
	"errors"
	"fmt"
)

type (
	LogLevel int
	Logger   struct {
		logLevel LogLevel
	}
)

const (
	Error LogLevel = iota
	Warn
	Info
	Debug
)

var logLevelIndex = map[string]LogLevel{
	"error": Error,
	"warn":  Warn,
	"info":  Info,
	"debug": Debug,
}

func New(level string) (*Logger, error) {
	inputLvl, ok := logLevelIndex[level]
	if !ok {
		return nil, errors.New("unknown log level: " + level)
	}
	return &Logger{logLevel: inputLvl}, nil
}

func (l Logger) Info(msg string) {
	if l.logLevel >= Info {
		fmt.Println(msg)
	}
}

func (l Logger) Debug(msg string) {
	if l.logLevel >= Debug {
		fmt.Println(msg)
	}
}

func (l Logger) Warn(msg string) {
	if l.logLevel >= Warn {
		fmt.Println(msg)
	}
}

func (l Logger) Error(msg string) {
	if l.logLevel >= Error {
		fmt.Println(msg)
	}
}
