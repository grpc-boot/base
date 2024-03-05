package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	defaultTickSecond = 5
	defaultMaxDays    = 7
)

var (
	defaultEncoder = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey: "Message",
		LevelKey:   "Level",
		TimeKey:    "DateTime",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		CallerKey:    "File",
		EncodeCaller: zapcore.ShortCallerEncoder,
	})

	timeFlag = func() string {
		return time.Now().Format("060102")
	}
)

var (
	zapLogger *Logger
)

func InitZapWithOption(conf Option, zapOpts ...zap.Option) error {
	logger, err := NewLogger(conf, zapOpts...)
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

func ZapDebug(msg string, fields ...zap.Field) {
	if zapLogger == nil {
		return
	}

	zapLogger.Debug(msg, fields...)
}

func ZapInfo(msg string, fields ...zap.Field) {
	if zapLogger == nil {
		return
	}

	zapLogger.Info(msg, fields...)
}

func ZapWarn(msg string, fields ...zap.Field) {
	if zapLogger == nil {
		return
	}

	zapLogger.Warn(msg, fields...)
}

func ZapError(msg string, fields ...zap.Field) {
	if zapLogger == nil {
		return
	}

	zapLogger.Error(msg, fields...)
}

func ZapDPanic(msg string, fields ...zap.Field) {
	if zapLogger == nil {
		return
	}

	zapLogger.DPanic(msg, fields...)
}

func ZapPanic(msg string, fields ...zap.Field) {
	if zapLogger == nil {
		return
	}

	zapLogger.Panic(msg, fields...)
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
