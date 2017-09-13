// Copyright 2017 Igor Dolzhikov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package standard

import (
	"log"
	"os"

	"github.com/takama/k8sapp/pkg/config"
	"github.com/takama/k8sapp/pkg/logger"
)

// UTC contains default UTC suffix
const UTC = "+0000 UTC "

// New returns logger that is compatible with the Logger interface
func New(cfg *logger.Config) logger.Logger {
	var flags int
	prefix := "[" + config.SERVICENAME + ":" + cfg.Level.String() + "] "
	if cfg.Out == nil {
		cfg.Out = os.Stdout
	}
	if cfg.Err == nil {
		cfg.Err = os.Stderr
	}
	if cfg.Time {
		flags = log.Ldate | log.Ltime | log.Lmicroseconds
		if cfg.UTC {
			flags = flags | log.LUTC
		}
	}
	return &stdLogger{
		Level:  cfg.Level,
		Time:   cfg.Time,
		UTC:    cfg.UTC,
		stdlog: log.New(cfg.Out, prefix, flags),
		errlog: log.New(cfg.Err, prefix, flags),
	}
}

// stdLogger implements the Logger interface
// except of using logger.Fields
type stdLogger struct {
	logger.Level
	Time   bool
	UTC    bool
	stdlog *log.Logger
	errlog *log.Logger
}

// Debug logs a debug message
func (l *stdLogger) Debug(v ...interface{}) {
	if l.Level == logger.LevelDebug {
		l.setStdPrefix(logger.LevelDebug)
		l.printStd(v...)
	}
}

// Debug logs a debug message with format
func (l *stdLogger) Debugf(format string, v ...interface{}) {
	if l.Level == logger.LevelDebug {
		l.setStdPrefix(logger.LevelDebug)
		l.printfStd(format, v...)
	}
}

// Info logs a info message
func (l *stdLogger) Info(v ...interface{}) {
	if l.Level <= logger.LevelInfo {
		l.setStdPrefix(logger.LevelInfo)
		l.printStd(v...)
	}
}

// Info logs a info message with format
func (l *stdLogger) Infof(format string, v ...interface{}) {
	if l.Level <= logger.LevelInfo {
		l.setStdPrefix(logger.LevelInfo)
		l.printfStd(format, v...)
	}
}

// Warn logs a warning message.
func (l *stdLogger) Warn(v ...interface{}) {
	if l.Level <= logger.LevelWarn {
		l.setStdPrefix(logger.LevelWarn)
		l.printStd(v...)
	}
}

// Warn logs a warning message with format.
func (l *stdLogger) Warnf(format string, v ...interface{}) {
	if l.Level <= logger.LevelWarn {
		l.setStdPrefix(logger.LevelWarn)
		l.printfStd(format, v...)
	}
}

// Error logs an error message
func (l *stdLogger) Error(v ...interface{}) {
	if l.Level <= logger.LevelError {
		l.setErrPrefix(logger.LevelError)
		l.printErr(v...)
	}
}

// Error logs an error message with format
func (l *stdLogger) Errorf(format string, v ...interface{}) {
	if l.Level <= logger.LevelError {
		l.setErrPrefix(logger.LevelError)
		l.printfErr(format, v...)
	}
}

// Fatal logs an error message followed by a call to os.Exit(1)
func (l *stdLogger) Fatal(v ...interface{}) {
	if l.Level <= logger.LevelFatal {
		l.setErrPrefix(logger.LevelFatal)
		l.printErr(v...)
	}
}

// Fatalf logs an error message with format followed by a call to ox.Exit(1)
func (l *stdLogger) Fatalf(format string, v ...interface{}) {
	if l.Level <= logger.LevelFatal {
		l.setErrPrefix(logger.LevelFatal)
		l.printfErr(format, v...)
	}
}

func (l *stdLogger) printStd(v ...interface{}) {
	if l.Time && l.UTC {
		l.stdlog.Print(append([]interface{}{UTC}, v...)...)
	} else {
		l.stdlog.Print(v...)
	}
}

func (l *stdLogger) printfStd(format string, v ...interface{}) {
	if l.Time && l.UTC {
		l.stdlog.Printf("%s"+format, append([]interface{}{UTC}, v...)...)
	} else {
		l.stdlog.Printf(format, v...)
	}
}

func (l *stdLogger) printErr(v ...interface{}) {
	if l.Time && l.UTC {
		l.errlog.Print(append([]interface{}{UTC}, v...)...)
	} else {
		l.errlog.Print(v...)
	}
}

func (l *stdLogger) printfErr(format string, v ...interface{}) {
	if l.Time && l.UTC {
		l.errlog.Printf("%s"+format, append([]interface{}{UTC}, v...)...)
	} else {
		l.errlog.Printf(format, v...)
	}
}

func (l *stdLogger) setStdPrefix(level logger.Level) {
	l.stdlog.SetPrefix("[" + config.SERVICENAME + ":" + level.String() + "] ")
}

func (l *stdLogger) setErrPrefix(level logger.Level) {
	l.errlog.SetPrefix("[" + config.SERVICENAME + ":" + level.String() + "] ")
}
