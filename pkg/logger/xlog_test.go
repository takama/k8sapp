package logger

import (
	"os"
	"testing"
)

func TestNewXLog(t *testing.T) {
	log1 := newXLog(&Config{
		Level: LevelDebug,
	})
	if log1 == nil {
		t.Error("Got uninitialized XLog logger")
	}
	log2 := newXLog(&Config{
		Level: LevelInfo,
		Out:   os.Stdout,
		Err:   os.Stdout,
	})
	if log2 == nil {
		t.Error("Got uninitialized XLog logger")
	}
}
