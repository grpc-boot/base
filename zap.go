package base

import (
	"github.com/grpc-boot/base/core/zaplogger"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	zapLogger *zaplogger.Logger
)

func InitZapWithOption(conf zaplogger.Option, zapOpts ...zap.Option) error {
	logger, err := zaplogger.NewLogger(conf, zapOpts...)
	if err != nil {
		return err
	}

	zapLogger = logger
	return nil
}

func IsLevel(level zapcore.Level) bool {
	if zapLogger == nil {
		return false
	}

	return zapLogger.Is(level)
}

func Debug(msg string, fields ...zap.Field) {
	ZapDebug(msg, fields...)
}

func ZapDebug(msg string, fields ...zap.Field) {
	if zapLogger == nil {
		return
	}

	zapLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	ZapInfo(msg, fields...)
}

func ZapInfo(msg string, fields ...zap.Field) {
	if zapLogger == nil {
		return
	}

	zapLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	ZapWarn(msg, fields...)
}

func ZapWarn(msg string, fields ...zap.Field) {
	if zapLogger == nil {
		return
	}

	zapLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	ZapError(msg, fields...)
}

func ZapError(msg string, fields ...zap.Field) {
	if zapLogger == nil {
		return
	}

	zapLogger.Error(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
	ZapDPanic(msg, fields...)
}

func ZapDPanic(msg string, fields ...zap.Field) {
	if zapLogger == nil {
		return
	}

	zapLogger.DPanic(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	ZapPanic(msg, fields...)
}

func ZapPanic(msg string, fields ...zap.Field) {
	if zapLogger == nil {
		return
	}

	zapLogger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	ZapFatal(msg, fields...)
}

func ZapFatal(msg string, fields ...zap.Field) {
	if zapLogger == nil {
		return
	}

	zapLogger.Fatal(msg, fields...)
}

func ZapSync() error {
	if zapLogger == nil {
		return nil
	}

	return zapLogger.Sync()
}
