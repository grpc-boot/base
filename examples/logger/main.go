package main

import (
	"os"
	"time"

	"github.com/grpc-boot/base"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	formal()
	log()

	level()
	log()
}

func log() {
	base.Debug("aa", "bb", "cc", 3, 5, 788.89)
	base.Info("aa", "bb", "cc", 3, 5, 788.89)
	base.Infof("msg %s %d", "as", 34)

	base.ZapDebug("run main",
		zap.Int64("datetime", time.Now().Unix()),
		zap.String("filename", "logger"),
	)

	base.ZapInfo("run main",
		zap.Int64("datetime", time.Now().Unix()),
		zap.String("filename", "logger"),
	)

	base.ZapWarn("please warn", zap.Int64("day", 15))

	base.ZapError("run main",
		zap.Int64("datetime", time.Now().Unix()),
		zap.String("filename", "logger"),
	)

	base.ZapSync()
}

func formal() {
	conf := zap.Config{
		Level:    zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "Message",
			LevelKey:   "Level",
			TimeKey:    "DateTime",
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format("2006-01-02 15:04:05"))
			},
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			CallerKey:    "File",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout", "info.log"},
		ErrorOutputPaths: []string{"stderr", "error.log"},
		InitialFields:    map[string]interface{}{"App": "boot"},
	}

	err := base.InitZapWithConfig(conf)

	if err != nil {
		base.RedFatal("init zap logger err:%s", err.Error())
	}
}

func level() {
	conf1 := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
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

	infoLog, _ := os.Create("info.log")
	errorLog, _ := os.Create("error.log")

	core := zapcore.NewTee(
		zapcore.NewCore(conf1, infoLog, zap.LevelEnablerFunc(func(z zapcore.Level) bool {
			return z >= zap.InfoLevel && z <= zap.WarnLevel
		})),
		zapcore.NewCore(conf1, errorLog, zap.LevelEnablerFunc(func(z zapcore.Level) bool {
			return z >= zap.ErrorLevel
		})),
	)

	base.InitZapWithCore(core)
}
