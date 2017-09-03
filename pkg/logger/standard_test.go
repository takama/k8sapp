package logger

import (
	"bytes"
	"strings"
	"testing"

	"github.com/takama/k8sapp/pkg/config"
)

func logMessage(level Level, message string, out, err *bytes.Buffer, time, utc bool) {
	log := New(&Config{
		Level: LevelDebug,
		Out:   out,
		Err:   err,
		Time:  time,
		UTC:   utc,
	})
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
}

func logMessagef(level Level, format, message string, out, err *bytes.Buffer, time, utc bool) {
	log := New(&Config{
		Level: LevelDebug,
		Out:   out,
		Err:   err,
		Time:  time,
		UTC:   utc,
	})
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

func testOutput(t *testing.T, level Level, message string, formated bool) {
	var want string
	prefix := "[" + config.SERVICENAME + ":" + level.String() + "] "
	out := &bytes.Buffer{}
	err := &bytes.Buffer{}
	if formated {
		want = prefix + message + "\n"
		logMessage(level, message, out, err, false, false)
	} else {
		want = prefix + "message=" + message + "\n"
		format := "message=%s"
		logMessagef(level, format, message, out, err, false, false)
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

func checkEmptyMessage(t *testing.T, out *bytes.Buffer, messageLevel, outputlevel Level) {
	if out.String() == "" {
		t.Errorf("Got empty %s message for %s output level", messageLevel.String(), outputlevel.String())
	}

}

func checkNonEmptyMessage(t *testing.T, out *bytes.Buffer, messageLevel, outputlevel Level) {
	if out.String() != "" {
		t.Errorf("Got non-empty %s message for %s output level", messageLevel.String(), outputlevel.String())
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
			checkEmptyMessage(t, out, messageLevel, level)
		default:
			checkNonEmptyMessage(t, out, messageLevel, level)
		}
	case LevelInfo:
		log.Info(message)
		switch level {
		case LevelDebug, LevelInfo:
			checkEmptyMessage(t, out, messageLevel, level)
		default:
			checkNonEmptyMessage(t, out, messageLevel, level)
		}
	case LevelWarn:
		log.Warn(message)
		switch level {
		case LevelDebug, LevelInfo, LevelWarn:
			checkEmptyMessage(t, out, messageLevel, level)
		default:
			checkNonEmptyMessage(t, out, messageLevel, level)
		}
	case LevelError:
		log.Error(message)
		switch level {
		case LevelDebug, LevelInfo, LevelWarn, LevelError:
			checkEmptyMessage(t, err, messageLevel, level)
		default:
			checkNonEmptyMessage(t, err, messageLevel, level)
		}
	case LevelFatal:
		log.Fatal(message)
		checkEmptyMessage(t, err, messageLevel, level)
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
	prefix := "[" + config.SERVICENAME + ":" + level.String() + "] "
	want := prefix + "__TIME__ " + UTC + message + "\n"
	out := &bytes.Buffer{}
	err := &bytes.Buffer{}
	logMessage(level, message, out, err, true, true)
	if level == LevelDebug || level == LevelInfo || level == LevelWarn {
		if got := out.String(); !strings.Contains(got, UTC) ||
			!strings.Contains(got, prefix) || !strings.Contains(got, message) {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	} else {
		if got := err.String(); !strings.Contains(got, UTC) ||
			!strings.Contains(got, prefix) || !strings.Contains(got, message) {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	}
}

func TestLogWithTime(t *testing.T) {
	for _, level := range []Level{LevelDebug, LevelInfo, LevelWarn, LevelError, LevelFatal} {
		testOutputWithTime(t, level, level.String()+" message")
	}
}
