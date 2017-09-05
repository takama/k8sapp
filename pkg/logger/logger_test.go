package logger

import "testing"

const (
	customLevel       Level = 17
	customLevelString       = "17"
)

func TestLevelString(t *testing.T) {
	for _, l := range []Level{LevelDebug, LevelInfo, LevelWarn, LevelError, LevelFatal} {
		if len(l.String()) <= 1 {
			t.Errorf("Got empty string for %v log level", l)
		}
	}
	if customLevel.String() != customLevelString {
		t.Error("All undefined levels should be presented as string")
		t.Errorf("invalid log level:\ngot:  %s\nwant: %s", customLevel.String(), customLevelString)
	}
}
