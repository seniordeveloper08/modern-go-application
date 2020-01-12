package commonadapter

import (
	"context"
	"fmt"
	"testing"

	"logur.dev/logur"
	"logur.dev/logur/logtesting"
)

func TestLogger_Levels(t *testing.T) {
	tests := map[string]struct {
		logFunc func(logger *Logger, msg string, fields ...map[string]interface{})
	}{
		"trace": {
			logFunc: (*Logger).Trace,
		},
		"debug": {
			logFunc: (*Logger).Debug,
		},
		"info": {
			logFunc: (*Logger).Info,
		},
		"warn": {
			logFunc: (*Logger).Warn,
		},
		"error": {
			logFunc: (*Logger).Error,
		},
	}

	for name, test := range tests {
		name, test := name, test

		t.Run(name, func(t *testing.T) {
			testLogger := logur.NewTestLogger()
			logger := NewLogger(testLogger)

			test.logFunc(logger, fmt.Sprintf("message: %s", name))

			level, _ := logur.ParseLevel(name)

			event := logur.LogEvent{
				Level: level,
				Line:  "message: " + name,
			}

			logtesting.AssertLogEventsEqual(t, event, *(testLogger.LastEvent()))
		})
	}
}

func TestLogger_Levels_Context(t *testing.T) {
	tests := map[string]struct {
		logFunc func(logger *Logger, ctx context.Context, msg string, fields ...map[string]interface{})
	}{
		"trace": {
			logFunc: (*Logger).TraceContext,
		},
		"debug": {
			logFunc: (*Logger).DebugContext,
		},
		"info": {
			logFunc: (*Logger).InfoContext,
		},
		"warn": {
			logFunc: (*Logger).WarnContext,
		},
		"error": {
			logFunc: (*Logger).ErrorContext,
		},
	}

	for name, test := range tests {
		name, test := name, test

		t.Run(name, func(t *testing.T) {
			testLogger := logur.NewTestLogger()
			logger := NewLogger(testLogger)

			test.logFunc(logger, context.Background(), fmt.Sprintf("message: %s", name))

			level, _ := logur.ParseLevel(name)

			event := logur.LogEvent{
				Level: level,
				Line:  "message: " + name,
			}

			logtesting.AssertLogEventsEqual(t, event, *(testLogger.LastEvent()))
		})
	}
}

func TestLogger_WithFields(t *testing.T) {
	testLogger := logur.NewTestLogger()

	fields := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}

	logger := NewLogger(testLogger).WithFields(fields)

	logger.Debug("message", nil)

	event := logur.LogEvent{
		Level:  logur.Debug,
		Line:   "message",
		Fields: fields,
	}

	logtesting.AssertLogEventsEqual(t, event, *(testLogger.LastEvent()))
}

type contextExtractor struct{}

func (*contextExtractor) Extract(ctx context.Context) map[string]interface{} {
	return map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}
}

func TestLogger_WithContext(t *testing.T) {
	testLogger := logur.NewTestLogger()

	logger := NewLogger(testLogger).WithContext(context.Background())

	logger.Debug("message", nil)

	event := logur.LogEvent{
		Level: logur.Debug,
		Line:  "message",
	}

	logtesting.AssertLogEventsEqual(t, event, *(testLogger.LastEvent()))
}

func TestContextAwareLogger_WithContext(t *testing.T) {
	testLogger := logur.NewTestLogger()

	logger := NewContextAwareLogger(testLogger, &contextExtractor{}).WithContext(context.Background())

	logger.Debug("message", nil)

	event := logur.LogEvent{
		Level: logur.Debug,
		Line:  "message",
		Fields: map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		},
	}

	logtesting.AssertLogEventsEqual(t, event, *(testLogger.LastEvent()))
}
