package zaplogger

import (
	"time"

	"go.uber.org/zap/zapcore"
)

var (
	DefaultEncoder = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
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

	TimeFlag = func() string {
		return time.Now().Format("06-01-02")
	}
)
