package logger

import "testing"

func TestLogrusLevel(t *testing.T) {
	for _, l := range []Level{LevelDebug, LevelInfo, LevelWarn, LevelError, LevelFatal} {
		if logrusLevelConverter(l) == 0 {
			t.Errorf("Got empty data for %s log level", l.String())
		}
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
