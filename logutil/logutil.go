// Package logutil has utilities for working with the Go Kit log package.
package logutil

import (
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// New will create a base logger with some commonly used / consistent settings
func New() log.Logger {
	logger := log.NewJSONLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	return logger
}

// Fatal logs a error message and exits the process.
func Fatal(logger log.Logger, args ...interface{}) {
	level.Info(logger).Log(args...)
	os.Exit(1)
}

// SetLevelKey changes the "level" key in a Go Kit logger, allowing the user
// to set it to something else. Useful for deploying services to GCP, as
// stackdriver expects a "severity" key instead.
//
// see https://github.com/go-kit/kit/issues/503
func SetLevelKey(logger log.Logger, key interface{}) log.Logger {
	return log.LoggerFunc(func(keyvals ...interface{}) error {
		for i := 1; i < len(keyvals); i += 2 {
			if _, ok := keyvals[i].(level.Value); ok {
				// overwriting the key without copying keyvals
				// techically violates the log.Logger contract
				// but is safe in this context because none
				// of the loggers in this program retain a reference
				// to keyvals
				keyvals[i-1] = key
				break
			}
		}
		return logger.Log(keyvals...)
	})
}
