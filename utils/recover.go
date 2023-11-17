package utils

import (
	"github.com/grpc-boot/base/v2/logger"
)

func Recover(eventName string, f func()) {
	defer func() {
		if err := recover(); err != nil {
			logger.ZapError("panic error",
				logger.Error(err.(error)),
				logger.Event(eventName),
			)
		}
	}()

	f()
}
