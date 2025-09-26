package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewZapLoggerWithLumberjack(basePath string, enc zapcore.Encoder, ljOption *lumberjack.Logger) *zap.Logger {
	basePath = strings.TrimSuffix(basePath, "/")
	cores := []zapcore.Core{
		zapcore.NewCore(enc, loadWriterSync(basePath+"/debug.log", ljOption), zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l == zapcore.DebugLevel
		})),
		zapcore.NewCore(enc, loadWriterSync(basePath+"/info.log", ljOption), zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l == zapcore.InfoLevel
		})),
		zapcore.NewCore(enc, loadWriterSync(basePath+"/warn.log", ljOption), zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l == zapcore.WarnLevel
		})),
		zapcore.NewCore(enc, loadWriterSync(basePath+"/error.log", ljOption), zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l == zapcore.ErrorLevel
		})),
	}

	core := zapcore.NewTee(cores...)
	return zap.New(core, zap.AddCaller())
}

func loadWriterSync(fileName string, opt *lumberjack.Logger) zapcore.WriteSyncer {
	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    opt.MaxSize,
		MaxAge:     opt.MaxAge,
		MaxBackups: opt.MaxBackups,
		LocalTime:  opt.LocalTime,
		Compress:   opt.Compress,
	})
}
