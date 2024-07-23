package logger

import (
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/grpc-boot/base/v2/internal"

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

	go logger.clearDirty()

	if _, err = logger.logger(); err != nil {
		return nil, err
	}

	go logger.checkFile()

	return logger, nil
}

func (l *Logger) checkFile() {
	if l.opt.TickSecond < 0 || l.opt.Path == "" {
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

			if err = l.clearDirty(); err != nil {
				fmt.Printf("clear dirty file error:%s\n", err.Error())
			}
		}()
	}
}

func (l *Logger) flagFileName(level, flag string) string {
	return fmt.Sprintf("%s/%s-%s.log", strings.TrimSuffix(l.opt.Path, "/"), level, flag)
}

func (l *Logger) flagFile(level, flag string) (*os.File, error) {
	fileName := l.flagFileName(level, flag)
	_ = internal.MkDir(path.Dir(fileName), 0766)
	return os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
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
	if l.opt.Path == "" {
		ow := zapcore.Lock(os.Stdout)
		ew := zapcore.Lock(os.Stderr)
		for ll := level; ll <= zapcore.FatalLevel; ll++ {
			if ll < zapcore.ErrorLevel {
				coreList = append(coreList, zapcore.NewCore(l.opt.encoder, ow, l.levelEnabler(ll)))
				continue
			}
			coreList = append(coreList, zapcore.NewCore(l.opt.encoder, ew, l.levelEnabler(ll)))
		}
	} else {
		for ll := level; ll <= zapcore.FatalLevel; ll++ {
			file, err := l.flagFile(ll.String(), flag)
			if err != nil {
				return nil, err
			}

			coreList = append(coreList, zapcore.NewCore(l.opt.encoder, file, l.levelEnabler(ll)))
		}
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

func (l *Logger) clearDirty() (err error) {
	entryList, err := os.ReadDir(l.opt.Path)
	if err != nil {
		return
	}

	now := time.Now()
	for _, entry := range entryList {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()

		if !strings.HasSuffix(fileName, ".log") {
			continue
		}

		end := strings.Index(fileName, "-")
		if end < 4 {
			continue
		}

		switch fileName[:end] {
		case zapcore.DebugLevel.String():
		case zapcore.InfoLevel.String():
		case zapcore.WarnLevel.String():
		case zapcore.ErrorLevel.String():
		case zapcore.DPanicLevel.String():
		case zapcore.PanicLevel.String():
		case zapcore.FatalLevel.String():
		default:
			continue
		}

		fullFileName := strings.TrimSuffix(l.opt.Path, "/") + "/" + fileName
		ctime, _, _, _ := internal.FileTime(fullFileName)
		if ctime < 1 || now.Unix()-ctime < int64(l.opt.MaxDays)*3600*24 {
			continue
		}

		_ = os.Remove(fullFileName)
	}
	return
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
//
//	entries. Applications should take care to call Sync before exiting.
func (l *Logger) Sync() error {
	logger, err := l.logger()
	if err != nil {
		return err
	}

	return logger.Sync()
}
