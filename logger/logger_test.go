package logger

import (
	"math/rand"
	"testing"

	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
)

func TestNewZapLogger(t *testing.T) {
	logger := NewZapLoggerWithLumberjack("./", DefaultJsonEncoder, &lumberjack.Logger{
		MaxSize:    10,
		MaxAge:     2,
		MaxBackups: 0,
		LocalTime:  true,
	})

	for i := 0; i < 10; i++ {
		switch rand.Uint32() % 4 {
		case 0:
			logger.Info("info", zap.String("key", "value"))
		case 1:
			logger.Warn("warn", zap.String("key", "value"))
		case 2:
			logger.Error("error", zap.String("key", "value"))
		case 3:
			logger.Debug("debug", zap.String("key", "value"))
		}
	}
}
