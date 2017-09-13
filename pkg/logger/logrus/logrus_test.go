package logrus

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/takama/k8sapp/pkg/logger"
)

const (
	customLevel logger.Level = 17
)

func TestLogrusLevel(t *testing.T) {
	for _, l := range []logger.Level{
		logger.LevelDebug,
		logger.LevelInfo,
		logger.LevelWarn,
		logger.LevelError,
		logger.LevelFatal,
	} {
		if logrusLevelConverter(l) == 0 {
			t.Errorf("Got empty data for %s log level", l.String())
		}
	}
	level := logrusLevelConverter(customLevel)
	if level != logrus.InfoLevel {
		t.Errorf("invalid log level:\ngot:  %s\nwant: %s", level, logrus.InfoLevel)
	}
}

func TestNewLogrus(t *testing.T) {
	log := New(&logger.Config{
		Level: logger.LevelDebug,
	})
	if log == nil {
		t.Error("Got uninitialized logrus logger")
	}
}
