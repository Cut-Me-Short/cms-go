/*
CutMeShort Go SDK - Simple Logger

Package logger provides basic structured logging for the SDK.
*/

package logger

import (
	"fmt"
	"log"
	"os"
)

// SimpleLogger provides basic structured logging functionality
type SimpleLogger struct {
	debugMode bool
	name      string
}

// New creates a new logger instance
func New(name string, debugMode bool) *SimpleLogger {
	return &SimpleLogger{
		name:      name,
		debugMode: debugMode,
	}
}

// Debug logs a debug message (only if debug mode is enabled)
func (l *SimpleLogger) Debug(msg string, keysAndValues ...interface{}) {
	if l.debugMode {
		l.log("DEBUG", msg, keysAndValues...)
	}
}

// Info logs an info message
func (l *SimpleLogger) Info(msg string, keysAndValues ...interface{}) {
	l.log("INFO", msg, keysAndValues...)
}

// Warn logs a warning message
func (l *SimpleLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.log("WARN", msg, keysAndValues...)
}

// Error logs an error message
func (l *SimpleLogger) Error(msg string, keysAndValues ...interface{}) {
	l.log("ERROR", msg, keysAndValues...)
}

// log internal logging function
func (l *SimpleLogger) log(level string, msg string, keysAndValues ...interface{}) {
	var fields string
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			if fields != "" {
				fields += ", "
			}
			fields += fmt.Sprintf("%v=%v", keysAndValues[i], keysAndValues[i+1])
		}
	}

	if fields != "" {
		log.Printf("[%s] [%s] %s { %s }", l.name, level, msg, fields)
	} else {
		log.Printf("[%s] [%s] %s", l.name, level, msg)
	}
}

// NoOpLogger implements the Logger interface but does nothing
type NoOpLogger struct{}

// Debug logs nothing
func (l *NoOpLogger) Debug(msg string, keysAndValues ...interface{}) {}

// Info logs nothing
func (l *NoOpLogger) Info(msg string, keysAndValues ...interface{}) {}

// Warn logs nothing
func (l *NoOpLogger) Warn(msg string, keysAndValues ...interface{}) {}

// Error logs nothing
func (l *NoOpLogger) Error(msg string, keysAndValues ...interface{}) {}

// NewNoOpLogger creates a no-op logger
func NewNoOpLogger() *NoOpLogger {
	return &NoOpLogger{}
}
