package zaplogger

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"go.uber.org/atomic"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	opt       *Option
	zapLogger atomic.Value
	mutex     sync.Mutex
	flag      string
	optList   []zap.Option
}

func NewLogger(opt Option, opts ...zap.Option) (logger *Logger, err error) {
	option := loadOption(opt)

	logger = &Logger{
		opt:     option,
		optList: opts,
	}

	if _, err = logger.logger(); err != nil {
		return nil, err
	}

	return logger, nil
}

func (l *Logger) flagFile(level, flag string) (*os.File, error) {
	filePath := fmt.Sprintf("%s/%s-%s.log", strings.TrimSuffix(l.opt.Path, "/"), level, flag)

	return os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
}

func (l *Logger) tee(flag string) (zapcore.Core, error) {
	debugFile, err := l.flagFile("debug", flag)
	if err != nil {
		return nil, err
	}

	infoFile, err := l.flagFile("info", flag)
	if err != nil {
		return nil, err
	}

	errorFile, err := l.flagFile("error", flag)
	if err != nil {
		return nil, err
	}

	level := zapcore.Level(l.opt.Level)

	tee := zapcore.NewTee(
		zapcore.NewCore(l.opt.encoder, debugFile, zap.LevelEnablerFunc(func(z zapcore.Level) bool {
			return z >= level && z >= zap.DebugLevel && z < zap.InfoLevel
		})),
		zapcore.NewCore(l.opt.encoder, infoFile, zap.LevelEnablerFunc(func(z zapcore.Level) bool {
			return z >= level && z >= zap.InfoLevel && z < zap.WarnLevel
		})),
		zapcore.NewCore(l.opt.encoder, errorFile, zap.LevelEnablerFunc(func(z zapcore.Level) bool {
			return z >= level && z >= zap.WarnLevel
		})),
	)

	return tee, nil
}

func (l *Logger) logger() (*zap.Logger, error) {
	flagStr := l.opt.flagFunc()

	if l.flag == flagStr {
		return l.zapLogger.Load().(*zap.Logger), nil
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.flag == flagStr {
		return l.zapLogger.Load().(*zap.Logger), nil
	}

	core, err := l.tee(flagStr)
	if err != nil {
		return nil, err
	}

	logger := zap.New(core, l.optList...)
	l.zapLogger.Store(logger)
	l.flag = flagStr

	return logger, nil
}

func (l *Logger) Is(level zapcore.Level) bool {
	return int8(level) >= l.opt.Level
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	if logger, _ := l.logger(); logger != nil {
		logger.Debug(msg, fields...)
	}
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	if logger, _ := l.logger(); logger != nil {
		logger.Info(msg, fields...)
	}
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	if logger, _ := l.logger(); logger != nil {
		logger.Warn(msg, fields...)
	}
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	if logger, _ := l.logger(); logger != nil {
		logger.Error(msg, fields...)
	}
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	if logger, _ := l.logger(); logger != nil {
		logger.Fatal(msg, fields...)
	}
}

func (l *Logger) Panic(msg string, fields ...zap.Field) {
	if logger, _ := l.logger(); logger != nil {
		logger.Panic(msg, fields...)
	}
}

func (l *Logger) Sync() error {
	logger, err := l.logger()
	if err != nil {
		return err
	}

	return logger.Sync()
}
