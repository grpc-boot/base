package mq

import "time"

var defaultOption = func() Option {
	return Option{
		ChanSize:            1024,
		RetryIntervalSecond: 3,
		MaxLength:           64 * 10000,
		MsgMinIdleSecond:    30,
	}
}

type Option struct {
	Group               string `json:"group" yaml:"group"`
	Consumer            string `json:"consumer" yaml:"consumer"`
	ChanSize            int64  `json:"chanSize" yaml:"chanSize"`
	RetryIntervalSecond int64  `json:"retryIntervalSecond" yaml:"retryIntervalSecond"`
	AutoCommit          bool   `json:"autoCommit" yaml:"autoCommit"`
	MaxLength           int64  `json:"maxLength" yaml:"maxLength"`
	MaxRetryTimes       int64  `json:"maxRetryTimes" yaml:"maxRetryTimes"`
	MsgMinIdleSecond    int64  `json:"msgMinIdleSecond" yaml:"msgMinIdleSecond"`
}

func formatOption(opt Option) Option {
	defaultOpt := defaultOption()
	defaultOpt.Group = opt.Group
	defaultOpt.Consumer = opt.Consumer
	defaultOpt.AutoCommit = opt.AutoCommit

	if opt.ChanSize > 0 {
		defaultOpt.ChanSize = opt.ChanSize
	}

	if opt.RetryIntervalSecond > 0 {
		defaultOpt.RetryIntervalSecond = opt.RetryIntervalSecond
	}

	if opt.MaxLength > 0 {
		defaultOpt.MaxLength = opt.MaxLength
	}

	if opt.MsgMinIdleSecond > 0 {
		defaultOpt.MsgMinIdleSecond = opt.MsgMinIdleSecond
	}

	return defaultOpt
}

func (o Option) RetryDuration() time.Duration {
	return time.Second * time.Duration(o.RetryIntervalSecond)
}

func (o Option) MsgMinIdle() time.Duration {
	return time.Second * time.Duration(o.MsgMinIdleSecond)
}
