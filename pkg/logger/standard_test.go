package logger

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/takama/k8sapp/pkg/config"
)

func testOutput(t *testing.T, level Level, message string, formated bool) {
	var want string
	prefix := "[" + config.SERVICENAME + ":" + level.String() + "] "
	out := &bytes.Buffer{}
	err := &bytes.Buffer{}
	log := New(&Config{
		Level: LevelDebug,
		Out:   out,
		Err:   err,
	})
	if formated {
		want = prefix + message + "\n"
		switch level {
		case LevelDebug:
			log.Debug(message)
		case LevelInfo:
			log.Info(message)
		case LevelWarn:
			log.Warn(message)
		case LevelError:
			log.Error(message)
		case LevelFatal:
			log.Fatal(message)
		}
	} else {
		want = prefix + "message=" + message + "\n"
		format := "message=%s"
		switch level {
		case LevelDebug:
			log.Debugf(format, message)
		case LevelInfo:
			log.Infof(format, message)
		case LevelWarn:
			log.Warnf(format, message)
		case LevelError:
			log.Errorf(format, message)
		case LevelFatal:
			log.Fatalf(format, message)
		}
	}
	if level == LevelDebug || level == LevelInfo || level == LevelWarn {
		if got := out.String(); got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	} else {
		if got := err.String(); got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	}
}

func TestLog(t *testing.T) {
	for _, level := range []Level{LevelDebug, LevelInfo, LevelWarn, LevelError, LevelFatal} {
		testOutput(t, level, level.String()+" message", false)
		testOutput(t, level, level.String()+" message", true)
	}
}

func testLevel(t *testing.T, level, messageLevel Level) {
	out := &bytes.Buffer{}
	err := &bytes.Buffer{}
	log := New(&Config{
		Level: level,
		Out:   out,
		Err:   err,
	})
	message := "message"
	switch messageLevel {
	case LevelDebug:
		log.Debug(message)
		switch level {
		case LevelDebug:
			if out.String() == "" {
				t.Errorf("Got empty debug message for %s output level", level.String())
			}
		default:
			if out.String() != "" {
				t.Errorf("Got non-empty debug message for %s output level", level.String())
			}
		}
	case LevelInfo:
		log.Info(message)
		switch level {
		case LevelDebug, LevelInfo:
			if out.String() == "" {
				t.Errorf("Got empty info message for %s output level", level.String())
			}
		default:
			if out.String() != "" {
				t.Errorf("Got non-empty info message for %s output level", level.String())
			}
		}
	case LevelWarn:
		log.Warn(message)
		switch level {
		case LevelDebug, LevelInfo, LevelWarn:
			if out.String() == "" {
				t.Errorf("Got empty warn message for %s output level", level.String())
			}
		default:
			if out.String() != "" {
				t.Errorf("Got non-empty warn message for %s output level", level.String())
			}
		}
	case LevelError:
		log.Error(message)
		switch level {
		case LevelDebug, LevelInfo, LevelWarn, LevelError:
			if err.String() == "" {
				t.Errorf("Got empty error message for %s output level", level.String())
			}
		default:
			if err.String() != "" {
				t.Errorf("Got non-empty error message for %s output level", level.String())
			}
		}
	case LevelFatal:
		log.Fatal(message)
		if err.String() == "" {
			t.Errorf("Got empty fatal message for %s output level", level.String())
		}
	}
}

func TestLevel(t *testing.T) {
	for _, level := range []Level{LevelDebug, LevelInfo, LevelWarn, LevelError, LevelFatal} {
		for _, messageLevel := range []Level{LevelDebug, LevelInfo, LevelWarn, LevelError, LevelFatal} {
			testLevel(t, level, messageLevel)
		}
	}
}

func testOutputWithTime(t *testing.T, level Level, message string) {
	var want string
	prefix := "[" + config.SERVICENAME + ":" + level.String() + "] "
	out := &bytes.Buffer{}
	err := &bytes.Buffer{}
	log := New(&Config{
		Level: LevelDebug,
		Out:   out,
		Err:   err,
		Time:  true,
		UTC:   true,
	})
	want = prefix + message + "\n"
	switch level {
	case LevelDebug:
		log.Debug(message)
	case LevelInfo:
		log.Info(message)
	case LevelWarn:
		log.Warn(message)
	case LevelError:
		log.Error(message)
	case LevelFatal:
		log.Fatal(message)
	}
	if level == LevelDebug || level == LevelInfo || level == LevelWarn {
		if got := out.String(); !strings.Contains(got, UTC) ||
			!strings.Contains(got, time.Now().Format("2006/01/02")) {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	} else {
		if got := err.String(); !strings.Contains(got, UTC) ||
			!strings.Contains(got, time.Now().Format("2006/01/02")) {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	}
}

func TestLogWithTime(t *testing.T) {
	for _, level := range []Level{LevelDebug, LevelInfo, LevelWarn, LevelError, LevelFatal} {
		testOutputWithTime(t, level, level.String()+" message")
	}
}
