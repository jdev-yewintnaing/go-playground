package logger

import (
	"fmt"
	"time"
)

type LogLevel string

const (
	InfoLevel  LogLevel = "INFO"
	ErrorLevel LogLevel = "ERROR"
	DebugLevel LogLevel = "DEBUG"
)

type LogWriter interface {
	WriteLog(message string) error
}

type ConsoleWriter struct {
}

func (cw ConsoleWriter) WriteLog(message string) error {
	fmt.Println(message)

	return nil
}

type Logger struct {
	dest   LogWriter
	prefix string
}

type LogError struct {
	Message   string
	Timestamp time.Time
	Err       error
}

func (e *LogError) Error() string {
	return fmt.Sprintf("[%s] %s: %v", e.Timestamp.Format(time.RFC3339), e.Message, e.Err)
}

func (l *Logger) Log(message string) error {
	err := l.dest.WriteLog(fmt.Sprintf("[%s] %s", l.prefix, message))

	if err != nil {
		return &LogError{
			Message:   err.Error(),
			Timestamp: time.Now(),
			Err:       err,
		}
	}

	return nil
}

type Options func(*Logger)

func WithPrefix(prefix string) Options {
	return func(l *Logger) {
		l.prefix = prefix
	}
}

func WithWriter(dest LogWriter) Options {
	return func(l *Logger) {
		l.dest = dest
	}
}

func NewLogger(opts ...Options) *Logger {
	logger := &Logger{
		dest:   ConsoleWriter{},
		prefix: "LOG",
	}

	for _, opt := range opts {
		opt(logger)
	}

	return logger
}

type MockWriter struct {
	LastMessage string
}

func (w *MockWriter) WriteLog(message string) error {
	w.LastMessage = message

	return nil
}
