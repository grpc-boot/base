package main

import (
	"time"

	"github.com/grpc-boot/base"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	conf := zap.Config{
		Level:    zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			EncodeLevel: zapcore.LowercaseLevelEncoder,
		},
		OutputPaths:      []string{"stdout", "boot.log"},
		ErrorOutputPaths: []string{"stderr", "error.log"},
		InitialFields:    map[string]interface{}{"app": "boot"},
	}

	err := base.InitWithConfig(conf)

	if err != nil {
		base.RedFatal("init zap logger err:%s", err.Error())
	}

	base.Debug("run main",
		zap.Int64("datetime", time.Now().Unix()),
		zap.String("filename", "logger"),
	)

	base.Info("run main",
		zap.Int64("datetime", time.Now().Unix()),
		zap.String("filename", "logger"),
	)

	base.Error("run main",
		zap.Int64("datetime", time.Now().Unix()),
		zap.String("filename", "logger"),
	)
}
