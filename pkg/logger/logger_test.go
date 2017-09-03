package logger

import "testing"

func TestLevelString(t *testing.T) {
	for _, l := range []Level{LevelDebug, LevelInfo, LevelWarn, LevelError, LevelFatal} {
		if len(l.String()) <= 1 {
			t.Errorf("Got empty string for %v log level", l)
		}
	}
}
