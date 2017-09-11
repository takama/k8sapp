package standard

import (
	"bytes"
	"strings"
	"testing"

	"github.com/takama/k8sapp/pkg/config"
	"github.com/takama/k8sapp/pkg/logger"
)

func TestNewLog(t *testing.T) {
	config := &logger.Config{}
	New(config)
	if config.Level != logger.LevelDebug {
		t.Errorf("Invalid level, got %s, want %s", config.Level, logger.LevelDebug)
	}
	if config.Out == nil {
		t.Error("Invalid logger output, got nil, want os.Stdout")
	}
	if config.Err == nil {
		t.Error("Invalid logger error output, got nil, want os.Stderr")
	}
}

func logMessage(level logger.Level, message string, out, err *bytes.Buffer, time, utc bool) {
	log := New(&logger.Config{
		Level: logger.LevelDebug,
		Out:   out,
		Err:   err,
		Time:  time,
		UTC:   utc,
	})
	switch level {
	case logger.LevelDebug:
		log.Debug(message)
	case logger.LevelInfo:
		log.Info(message)
	case logger.LevelWarn:
		log.Warn(message)
	case logger.LevelError:
		log.Error(message)
	case logger.LevelFatal:
		log.Fatal(message)
	}
}

func logMessageFormated(level logger.Level, format, message string, out, err *bytes.Buffer, time, utc bool) {
	log := New(&logger.Config{
		Level: logger.LevelDebug,
		Out:   out,
		Err:   err,
		Time:  time,
		UTC:   utc,
	})
	switch level {
	case logger.LevelDebug:
		log.Debugf(format, message)
	case logger.LevelInfo:
		log.Infof(format, message)
	case logger.LevelWarn:
		log.Warnf(format, message)
	case logger.LevelError:
		log.Errorf(format, message)
	case logger.LevelFatal:
		log.Fatalf(format, message)
	}
}

func testOutput(t *testing.T, level logger.Level, message string, formated bool) {
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
		logMessageFormated(level, format, message, out, err, false, false)
	}
	if level == logger.LevelDebug || level == logger.LevelInfo || level == logger.LevelWarn {
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
	for _, level := range []logger.Level{
		logger.LevelDebug,
		logger.LevelInfo,
		logger.LevelWarn,
		logger.LevelError,
		logger.LevelFatal,
	} {
		testOutput(t, level, level.String()+" message", false)
		testOutput(t, level, level.String()+" message", true)
	}
}

func checkEmptyMessage(t *testing.T, out *bytes.Buffer, messageLevel, outputlevel logger.Level) {
	if out.String() == "" {
		t.Errorf("Got empty %s message for %s output level", messageLevel, outputlevel)
	}

}

func checkNonEmptyMessage(t *testing.T, out *bytes.Buffer, messageLevel, outputlevel logger.Level) {
	if out.String() != "" {
		t.Errorf("Got non-empty %s message for %s output level", messageLevel, outputlevel)
	}
}

func testLevel(t *testing.T, level, messageLevel logger.Level) {
	out := &bytes.Buffer{}
	err := &bytes.Buffer{}
	log := New(&logger.Config{
		Level: level,
		Out:   out,
		Err:   err,
	})
	message := "message"
	switch messageLevel {
	case logger.LevelDebug:
		log.Debug(message)
		switch level {
		case logger.LevelDebug:
			checkEmptyMessage(t, out, messageLevel, level)
		default:
			checkNonEmptyMessage(t, out, messageLevel, level)
		}
	case logger.LevelInfo:
		log.Info(message)
		switch level {
		case logger.LevelDebug, logger.LevelInfo:
			checkEmptyMessage(t, out, messageLevel, level)
		default:
			checkNonEmptyMessage(t, out, messageLevel, level)
		}
	case logger.LevelWarn:
		log.Warn(message)
		switch level {
		case logger.LevelDebug, logger.LevelInfo, logger.LevelWarn:
			checkEmptyMessage(t, out, messageLevel, level)
		default:
			checkNonEmptyMessage(t, out, messageLevel, level)
		}
	case logger.LevelError:
		log.Error(message)
		switch level {
		case logger.LevelDebug, logger.LevelInfo, logger.LevelWarn, logger.LevelError:
			checkEmptyMessage(t, err, messageLevel, level)
		default:
			checkNonEmptyMessage(t, err, messageLevel, level)
		}
	case logger.LevelFatal:
		log.Fatal(message)
		checkEmptyMessage(t, err, messageLevel, level)
	}
}

func TestLevel(t *testing.T) {
	for _, level := range []logger.Level{
		logger.LevelDebug,
		logger.LevelInfo,
		logger.LevelWarn,
		logger.LevelError,
		logger.LevelFatal,
	} {
		for _, messageLevel := range []logger.Level{
			logger.LevelDebug,
			logger.LevelInfo,
			logger.LevelWarn,
			logger.LevelError,
			logger.LevelFatal,
		} {
			testLevel(t, level, messageLevel)
		}
	}
}

func testOutputWithTime(t *testing.T, level logger.Level, message string) {
	prefix := "[" + config.SERVICENAME + ":" + level.String() + "] "
	want := prefix + "__TIME__ " + UTC + message + "\n"
	out := &bytes.Buffer{}
	err := &bytes.Buffer{}
	logMessage(level, message, out, err, true, true)
	if level == logger.LevelDebug || level == logger.LevelInfo || level == logger.LevelWarn {
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

func testOutputFormatedWithTime(t *testing.T, level logger.Level, message string) {
	prefix := "[" + config.SERVICENAME + ":" + level.String() + "] "
	want := prefix + "__TIME__ " + UTC + message + "\n"
	out := &bytes.Buffer{}
	err := &bytes.Buffer{}
	logMessageFormated(level, "%s", message, out, err, true, true)
	if level == logger.LevelDebug || level == logger.LevelInfo || level == logger.LevelWarn {
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
	for _, level := range []logger.Level{
		logger.LevelDebug,
		logger.LevelInfo,
		logger.LevelWarn,
		logger.LevelError,
		logger.LevelFatal,
	} {
		testOutputWithTime(t, level, level.String()+" message")
		testOutputFormatedWithTime(t, level, level.String()+" message")
	}
}
