package utils

import (
	"fmt"

	"github.com/grpc-boot/base/v2/logger"
)

func Recover(eventName string, f func()) {
	defer func() {
		if err := recover(); err != nil {
			if _, ok := err.(error); !ok {
				err = fmt.Errorf("error:%v", err)
			}

			logger.ZapError("panic error",
				logger.Error(err.(error)),
				logger.Event(eventName),
			)
		}
	}()

	f()
}
