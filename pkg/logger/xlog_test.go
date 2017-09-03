package logger

import "testing"

func TestNewXLog(t *testing.T) {
	log := newXLog(&Config{
		Level: LevelDebug,
	})
	if log == nil {
		t.Error("Got uninitialized XLog logger")
	}
}
