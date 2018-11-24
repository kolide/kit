// +build windows

package logutil

import "github.com/go-kit/kit/log"

func swapLevelHandler(base log.Logger, swapLogger *log.SwapLogger, debug bool) {
	// noop for now
}
