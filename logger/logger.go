package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	DefaultLevel = zap.InfoLevel
	DefaultLog   = func(level zapcore.Level, msg string, fields ...zap.Field) {
		line, _ := EncodeFields(DefaultJsonEncoder, level, msg, fields...)
		if len(line) > 0 {
			fmt.Println(bytes2String(line))
		}
	}
)

func IsDebug() bool {
	return DefaultLevel.Enabled(zapcore.DebugLevel)
}

func Error(msg string, fields ...zap.Field) {
	if DefaultLevel.Enabled(zapcore.ErrorLevel) {
		DefaultLog(zap.ErrorLevel, msg, fields...)
	}
}

func Warn(msg string, fields ...zap.Field) {
	if DefaultLevel.Enabled(zapcore.WarnLevel) {
		DefaultLog(zap.WarnLevel, msg, fields...)
	}
}

func Info(msg string, fields ...zap.Field) {
	if DefaultLevel.Enabled(zapcore.InfoLevel) {
		DefaultLog(zap.InfoLevel, msg, fields...)
	}
}

func Debug(msg string, fields ...zap.Field) {
	if DefaultLevel.Enabled(zapcore.DebugLevel) {
		DefaultLog(zap.DebugLevel, msg, fields...)
	}
}
