package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"time"
)

type LogLevel string

const (
	LevelDebug LogLevel = "DEBUG"
	LevelInfo  LogLevel = "INFO"
	LevelError LogLevel = "ERROR"
)

type LogEntry struct {
	Level     LogLevel               `json:"level"`
	Message   string                 `json:"message"`
	Error     string                 `json:"error,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	Stack     string                 `json:"stack,omitempty"`
}

type Logger struct {
	output *json.Encoder
}

func New() *Logger {
	return &Logger{
		output: json.NewEncoder(os.Stdout),
	}
}

func (l *Logger) log(level LogLevel, msg string, err error) {
	entry := LogEntry{
		Level:     level,
		Message:   msg,
		Timestamp: time.Now(),
		Fields:    make(map[string]interface{}),
	}

	if err != nil {
		entry.Error = err.Error()
		entry.Stack = string(debug.Stack())
	}

	_, file, line, ok := runtime.Caller(2)
	if ok {
		entry.Fields["caller"] = fmt.Sprintf("%s:%d", file, line)
	}

	l.output.Encode(entry)
}

func (l *Logger) Debug(msg string) {
	l.log(LevelDebug, msg, nil)
}

func (l *Logger) Info(msg string) {
	l.log(LevelInfo, msg, nil)
}

func (l *Logger) Error(msg string, err error) {
	l.log(LevelError, msg, err)
}

func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	return &Logger{
		output: l.output,
	}
}
