package logger

import (
	"runtime"
	"time"

	"go.uber.org/zap/zapcore"
)

var (
	DefaultJsonEncoder = zapcore.NewJSONEncoder(
		zapcore.EncoderConfig{
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
	)
)

func EncodeFields(encoder zapcore.Encoder, level zapcore.Level, msg string, fields ...zapcore.Field) ([]byte, error) {
	ent := zapcore.Entry{
		Level:      level,
		Time:       time.Now(),
		LoggerName: "zap",
		Message:    msg,
	}

	_, file, line, ok := runtime.Caller(3)
	if ok {
		ent.Caller = zapcore.EntryCaller{
			Defined: true,
			File:    file,
			Line:    line,
		}
	}

	buf, err := encoder.EncodeEntry(ent, fields)
	if err != nil || buf == nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
