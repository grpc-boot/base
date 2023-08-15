package base

import "github.com/grpc-boot/base/core/zaplogger"

func Recover(eventName string, f func()) {
	defer func() {
		if err := recover(); err != nil {
			ZapError("panic error",
				zaplogger.Error(err.(error)),
				zaplogger.Event(eventName),
			)
		}
	}()

	f()
}
