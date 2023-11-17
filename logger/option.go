package logger

import "go.uber.org/zap/zapcore"

type Option struct {
	// 日志级别：debug: -1 info:0 warn: 1 error: 2 dpanic: 3 panic: 4 fatal: 5
	Level int8 `json:"level" yaml:"level"`
	// 日志目录，默认路径：/tmp
	Path string `json:"path" yaml:"path"`
	// 定时检测文件时间间隔，单位秒，配置小于0的值表示禁用，默认5秒
	TickSecond int32 `json:"tickSecond" yaml:"tickSecond"`

	encoder  zapcore.Encoder
	flagFunc func() string
}

func loadOption(o Option) *Option {
	if o.Path == "" {
		o.Path = defaultPath
	}

	if o.encoder == nil {
		o.encoder = defaultEncoder
	}

	if o.flagFunc == nil {
		o.flagFunc = timeFlag
	}

	if o.TickSecond == 0 {
		o.TickSecond = defaultTickSecond
	}

	return &o
}

type Options func(o *Option)

func (opt *Option) WithEncoder(encoder zapcore.Encoder) {
	opt.encoder = encoder
}

func (opt *Option) WithFlagFunc(flag func() string) {
	opt.flagFunc = flag
}
