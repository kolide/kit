// +build !windows

package logutil

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func swapLevelHandler(base log.Logger, swapLogger *log.SwapLogger, debug bool) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGUSR2)
	for {
		<-sigChan
		if debug {
			newLogger := level.NewFilter(base, level.AllowInfo())
			swapLogger.Swap(newLogger)
		} else {
			newLogger := level.NewFilter(base, level.AllowDebug())
			swapLogger.Swap(newLogger)
		}
		level.Info(swapLogger).Log("msg", "swapping level", "debug", !debug)
		debug = !debug
	}
}
