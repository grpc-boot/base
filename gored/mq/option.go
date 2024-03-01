package mq

import "time"

var defaultOption = func() Option {
	return Option{
		RetrySecond: 3,
		MaxLength:   64 * 10000,
	}
}

type Option struct {
	Group         string `json:"group" yaml:"group"`
	Consumer      string `json:"consumer" yaml:"consumer"`
	RetrySecond   int64  `json:"retrySecond" yaml:"retrySecond"`
	AutoCommit    bool   `json:"autoCommit" yaml:"autoCommit"`
	MaxLength     int64  `json:"maxLength" yaml:"maxLength"`
	MaxRetryCount int64  `json:"maxRetryCount" yaml:"maxRetryCount"`
}

func formatOption(opt Option) Option {
	defaultOpt := defaultOption()
	defaultOpt.Group = opt.Group
	defaultOpt.Consumer = opt.Consumer
	defaultOpt.AutoCommit = opt.AutoCommit

	if opt.RetrySecond > 0 {
		defaultOpt.RetrySecond = opt.RetrySecond
	}

	if opt.MaxLength > 0 {
		defaultOpt.MaxLength = opt.MaxLength
	}

	return defaultOpt
}

func (o Option) RetryDuration() time.Duration {
	return time.Second * time.Duration(o.RetrySecond)
}
