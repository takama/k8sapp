package logger

import (
	"github.com/rs/xlog"
)

// newXLog creates "github.com/rs/xlog" logger
func newXLog(config *Config) Logger {
	return xlog.New(xlog.Config{
		Level:  xlog.Level(config.Level),
		Fields: config.Fields,
		Output: xlog.NewConsoleOutput(),
	})
}
