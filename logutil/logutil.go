// Package logutil has utilities for working with the Go Kit log package.
package logutil

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

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

// NewServerLogger creates a standard logger for Kolide services.
// The logger will output JSON structured logs with a
// "severity" field set to either "info" or "debug".
// The acceptable level can be swapped by sending SIGUSR2 to the process.
func NewServerLogger(debug bool) log.Logger {
	base := log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	base = log.With(base, "ts", log.DefaultTimestampUTC)
	base = SetLevelKey(base, "severity")
	base = level.NewInjector(base, level.InfoValue())

	lev := level.AllowInfo()
	if debug {
		lev = level.AllowDebug()
	}

	base = log.With(base, "caller", log.Caller(6))

	var swapLogger log.SwapLogger
	swapLogger.Swap(level.NewFilter(base, lev))

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGUSR2)
	go func() {
		for {
			<-sigChan
			if debug {
				newLogger := level.NewFilter(base, level.AllowInfo())
				swapLogger.Swap(newLogger)
			} else {
				newLogger := level.NewFilter(base, level.AllowDebug())
				swapLogger.Swap(newLogger)
			}
			level.Info(&swapLogger).Log("msg", "swapping level", "debug", !debug)
			debug = !debug
		}
	}()
	return &swapLogger
}
