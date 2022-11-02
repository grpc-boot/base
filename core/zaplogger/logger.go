package zaplogger

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

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

	// 定时检测文件是否存在
	go logger.checkFile()

	return logger, nil
}

func (l *Logger) checkFile() {
	if l.opt.TickSecond < 0 {
		return
	}

	tick := time.NewTicker(time.Second * time.Duration(l.opt.TickSecond))
	for range tick.C {
		func() {
			defer func() {
				err := recover()
				if err != nil {
					fmt.Printf("check file panic error:%v\n", err)
				}
			}()

			var (
				level   = zapcore.Level(l.opt.Level)
				flagStr = l.opt.flagFunc()
			)

			var reset bool

			for ll := level; ll <= zapcore.FatalLevel; ll++ {
				fileName := l.flagFileName(ll.String(), flagStr)
				_, err := os.Stat(fileName)
				if os.IsNotExist(err) {
					reset = true
					break
				}
			}

			if !reset {
				return
			}

			l.mutex.Lock()
			defer l.mutex.Unlock()

			_, err := l.resetZap(flagStr)
			if err != nil {
				fmt.Printf("reset logger file error:%s\n", err.Error())
			}
		}()
	}
}

func (l *Logger) flagFileName(level, flag string) string {
	return fmt.Sprintf("%s/%s-%s.log", strings.TrimSuffix(l.opt.Path, "/"), level, flag)
}

func (l *Logger) flagFile(level, flag string) (*os.File, error) {
	return os.OpenFile(l.flagFileName(level, flag), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
}

func (l *Logger) levelEnabler(level zapcore.Level) zap.LevelEnablerFunc {
	nextLevel := level + 1
	return func(z zapcore.Level) bool {
		return z < nextLevel && z >= level
	}
}

func (l *Logger) tee(flag string) (zapcore.Core, error) {
	var (
		level = zapcore.Level(l.opt.Level)
	)

	if level > zapcore.FatalLevel || level < zapcore.DebugLevel {
		return nil, ErrInvalidLevel
	}

	coreList := make([]zapcore.Core, 0, zapcore.FatalLevel-level+1)
	for ll := level; ll <= zapcore.FatalLevel; ll++ {
		file, err := l.flagFile(ll.String(), flag)
		if err != nil {
			return nil, err
		}

		coreList = append(coreList, zapcore.NewCore(l.opt.encoder, file, l.levelEnabler(ll)))
	}

	tee := zapcore.NewTee(coreList...)
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

	return l.resetZap(flagStr)
}

func (l *Logger) resetZap(flagStr string) (*zap.Logger, error) {
	core, err := l.tee(flagStr)
	if err != nil {
		return nil, err
	}

	logger := zap.New(core, l.optList...)
	l.zapLogger.Store(logger)
	l.flag = flagStr

	return logger, nil
}

// Is 判断日志级别
func (l *Logger) Is(level zapcore.Level) bool {
	return int8(level) >= l.opt.Level
}

// Debug logs are typically voluminous, and are usually disabled in
// production.
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	if logger, _ := l.logger(); logger != nil {
		logger.Debug(msg, fields...)
	}
}

// Info is the default logging priority.
func (l *Logger) Info(msg string, fields ...zap.Field) {
	if logger, _ := l.logger(); logger != nil {
		logger.Info(msg, fields...)
	}
}

// Warn logs are more important than Info, but don't need individual
// human review.
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	if logger, _ := l.logger(); logger != nil {
		logger.Warn(msg, fields...)
	}
}

// Error logs are high-priority. If an application is running smoothly,
// it shouldn't generate any error-level logs.
func (l *Logger) Error(msg string, fields ...zap.Field) {
	if logger, _ := l.logger(); logger != nil {
		logger.Error(msg, fields...)
	}
}

// DPanic logs are particularly important errors. In development the
// logger panics after writing the message.
func (l *Logger) DPanic(msg string, fields ...zap.Field) {
	if logger, _ := l.logger(); logger != nil {
		logger.DPanic(msg, fields...)
	}
}

// Panic logs a message, then panics.
func (l *Logger) Panic(msg string, fields ...zap.Field) {
	if logger, _ := l.logger(); logger != nil {
		logger.Panic(msg, fields...)
	}
}

// Fatal logs a message, then calls os.Exit(1).
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	if logger, _ := l.logger(); logger != nil {
		logger.Fatal(msg, fields...)
	}
}

// Sync calls the underlying Core's Sync method, flushing any buffered log
//  entries. Applications should take care to call Sync before exiting.
func (l *Logger) Sync() error {
	logger, err := l.logger()
	if err != nil {
		return err
	}

	return logger.Sync()
}
