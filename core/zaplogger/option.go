package zaplogger

import "go.uber.org/zap/zapcore"

type Option struct {
	Level int8   `json:"level" yaml:"level"`
	Path  string `json:"path" yaml:"path"`

	encoder  zapcore.Encoder
	flagFunc func() string
}

func loadOption(o Option) *Option {
	if o.encoder == nil {
		o.encoder = DefaultEncoder
	}

	if o.flagFunc == nil {
		o.flagFunc = TimeFlag
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
