package base

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	zapLogger *zap.Logger
)

// InitZap 初始化Zap日志
func InitZap(l *zap.Logger) {
	InitLogger(l.Sugar())

	zapLogger = l
}

// InitZapWithConfig with zap.Config初始化Zap日志
func InitZapWithConfig(conf zap.Config, opts ...zap.Option) error {
	zLogger, err := conf.Build(opts...)
	if err != nil {
		return err
	}

	InitZap(zLogger)

	return nil
}

// InitZapWithCore with core初始化Zap日志
func InitZapWithCore(core zapcore.Core, opts ...zap.Option) {
	InitZap(zap.New(core, opts...))
}

// InitZapWithProduction with production初始化Zap日志
func InitZapWithProduction(opts ...zap.Option) error {
	zLogger, err := zap.NewProduction(opts...)
	if err != nil {
		return err
	}

	InitZap(zLogger)

	return nil
}

func ZapDebug(msg string, fields ...zap.Field) {
	zapLogger.Debug(msg, fields...)
}

func ZapInfo(msg string, fields ...zap.Field) {
	zapLogger.Info(msg, fields...)
}

func ZapWarn(msg string, fields ...zap.Field) {
	zapLogger.Warn(msg, fields...)
}

func ZapError(msg string, fields ...zap.Field) {
	zapLogger.Error(msg, fields...)
}

func ZapFatal(msg string, fields ...zap.Field) {
	zapLogger.Fatal(msg, fields...)
}

func ZapPanic(msg string, fields ...zap.Field) {
	zapLogger.Panic(msg, fields...)
}

func ZapSync() error {
	return zapLogger.Sync()
}
