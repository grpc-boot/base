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
	return zapLogger.Is(level)
}

func Debug(msg string, fields ...zap.Field) {
	zapLogger.Debug(msg, fields...)
}

func ZapDebug(msg string, fields ...zap.Field) {
	zapLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	zapLogger.Info(msg, fields...)
}

func ZapInfo(msg string, fields ...zap.Field) {
	zapLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	zapLogger.Warn(msg, fields...)
}

func ZapWarn(msg string, fields ...zap.Field) {
	zapLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	zapLogger.Error(msg, fields...)
}

func ZapError(msg string, fields ...zap.Field) {
	zapLogger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	zapLogger.Fatal(msg, fields...)
}

func ZapFatal(msg string, fields ...zap.Field) {
	zapLogger.Fatal(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	zapLogger.Panic(msg, fields...)
}

func ZapPanic(msg string, fields ...zap.Field) {
	zapLogger.Panic(msg, fields...)
}

func ZapSync() error {
	return zapLogger.Sync()
}
