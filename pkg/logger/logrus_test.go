package logger

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLogrusLevel(t *testing.T) {
	for _, l := range []Level{LevelDebug, LevelInfo, LevelWarn, LevelError, LevelFatal} {
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
	log := newLogrus(&Config{
		Level: LevelDebug,
	})
	if log == nil {
		t.Error("Got uninitialized logrus logger")
	}
}
