package utils

import (
	"github.com/grpc-boot/base/v2/logger"

	"go.uber.org/zap"
)

func Recover(traceId, action string, handler func()) {
	defer func() {
		err := recover()
		if err != nil {
			logger.Error("panic error",
				zap.String("Trace", traceId),
				zap.NamedError("Error", err.(error)),
				zap.String("Action", action),
				zap.Stack("Stack"),
			)
		}
	}()

	handler()
}
