package base

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
)

func Init(l *zap.Logger) {
	logger = l
}

func InitWithConfig(conf zap.Config, opts ...zap.Option) error {
	zapLogger, err := conf.Build(opts...)
	if err != nil {
		return err
	}

	Init(zapLogger)

	return nil
}

func InitWithCore(core zapcore.Core, opts ...zap.Option) {
	Init(zap.New(core, opts...))
}

func InitWithProduction(opts ...zap.Option) error {
	zapLogger, err := zap.NewProduction(opts...)
	if err != nil {
		return err
	}

	Init(zapLogger)

	return nil
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}

func Sync() error {
	return logger.Sync()
}
