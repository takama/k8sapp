package logger

import (
	"io"
	"strconv"
)

// Level defines log levels
type Level int

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	default:
		return strconv.Itoa(int(l))
	}
}

// Log levels
const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

// Fields represents a set of log message fields
type Fields map[string]interface{}

// Logger defines the interface for a compatible logger
type Logger interface {
	// Debug logs a debug message
	Debug(v ...interface{})
	// Debug logs a debug message with format
	Debugf(format string, v ...interface{})
	// Info logs a info message
	Info(v ...interface{})
	// Info logs a info message with format
	Infof(format string, v ...interface{})
	// Warn logs a warning message.
	Warn(v ...interface{})
	// Warn logs a warning message with format.
	Warnf(format string, v ...interface{})
	// Error logs an error message
	Error(v ...interface{})
	// Error logs an error message with format
	Errorf(format string, v ...interface{})
	// Fatal logs an error message followed by a call to os.Exit(1)
	Fatal(v ...interface{})
	// Fatalf logs an error message with format followed by a call to ox.Exit(1)
	Fatalf(format string, v ...interface{})
}

// Config contains log level and default fields
type Config struct {
	// Level is the maximum level to output, logs with lower level are discarded.
	Level
	// Fields defines default fields to use with all messages.
	Fields
	// Output destination for levels: debug, info, warn
	Out io.Writer
	// Output destination for levels: error, fatal
	Err io.Writer
	// Do not log date and time if false
	Time bool
	// Use UTC time
	UTC bool
}

// New returns new logger
func New(cfg *Config) Logger {
	// There should be any implementation which compatible with logger interface
	return newStdLog(cfg)
}
