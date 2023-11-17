package logger

import (
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewLogger(t *testing.T) {
	logger, err := NewLogger(Option{
		Level: int8(zapcore.InfoLevel),
		Path:  "./logs",
	})

	if err != nil {
		t.Fatalf("want nil, got %s\n", err.Error())
	}

	defer logger.Sync()

	logger.Debug("dddddd", zap.Int64("CurrentTime", time.Now().Unix()))
	logger.Info("iiiii", zap.Int64("CurrentTime", time.Now().Unix()))
	logger.Warn("wwwww", zap.Int64("CurrentTime", time.Now().Unix()))
	logger.Error("eeeee", zap.Int64("CurrentTime", time.Now().Unix()))
	logger.DPanic("ddddd", zap.Int64("CurrentTime", time.Now().Unix()))

	func() {
		defer func() {
			er := recover()
			if er != nil {
				t.Logf("panic msg: %s\n", er)
			}
		}()
		logger.Panic("ppppp", zap.Int64("CurrentTime", time.Now().Unix()))
	}()
}

func TestLogger_Info(t *testing.T) {
	logger, err := NewLogger(Option{
		Level: int8(zapcore.InfoLevel),
		Path:  "./logs",
	})

	if err != nil {
		t.Fatalf("want nil, got %s\n", err.Error())
	}

	defer logger.Sync()

	tick := time.NewTicker(time.Second)
	var num int
	for range tick.C {
		if num > 10 {
			return
		}
		logger.Info("iiiiiii", zap.Int64("Current", time.Now().Unix()))
		num++
	}
}

// go test -v -bench=. -benchtime=15s -run=BenchmarkLogger_Error
// BenchmarkLogger_Error-8          2794917              8421 ns/op
func BenchmarkLogger_Error(b *testing.B) {
	opt := Option{
		Level: int8(zapcore.InfoLevel),
		Path:  "./logs",
	}

	// 切割粒度为秒，默认切割粒度为日
	opt.WithFlagFunc(func() string {
		return time.Now().Format("06-01-02_15_04_05")
	})

	logger, err := NewLogger(opt)

	if err != nil {
		b.Fatalf("want nil, got %s\n", err.Error())
	}

	b.ResetTimer()

	defer logger.Sync()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Error("test with error",
				zap.Int64("UnixNano", time.Now().UnixNano()),
			)
		}
	})
}
