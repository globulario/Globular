// Package controlplane implements Envoy xDS control plane.
package controlplane

import (
	"log"
)

// Logger is a simple logger with support for debug mode.
type Logger struct {
	Debug bool
}

// Debugf logs a formatted debug message if debugging is enabled.
// The message is formatted according to the specified format and arguments.
// It appends a newline character to the message before logging.
func (logger Logger) Debugf(format string, args ...interface{}) {
	if logger.Debug {
		log.Printf(format+"\n", args...)
	}
}

// Infof logs an informational message formatted according to the specified format string and arguments.
// The message is only logged if the logger's Debug field is set to true.
// It appends a newline character to the formatted message before logging.
func (logger Logger) Infof(format string, args ...interface{}) {
	if logger.Debug {
		log.Printf(format+"\n", args...)
	}
}

// Warnf logs a warning message with formatting support.
// The message is formatted according to the specified format string and arguments.
// A newline character is appended to the formatted message.
func (logger Logger) Warnf(format string, args ...interface{}) {
	log.Printf(format+"\n", args...)
}

// Errorf logs an error message with formatting support.
// The message is formatted according to the specified format string and arguments.
// A newline character is appended to the formatted message before logging.
func (logger Logger) Errorf(format string, args ...interface{}) {

	log.Printf(format+"\n", args...)
}
