package queue

import "time"

var (
	DefaultOptions = func() Options {
		return Options{
			MaxRetry:        16,
			FetchCheckSec:   5,
			FetchTimeoutSec: 100,
			RetryCheckSec:   10,
			RetryTimeoutSec: 600,
		}
	}
)

type Options struct {
	MaxRetry        int64 `json:"maxRetry" yaml:"maxRetry"`
	FetchCheckSec   int64 `json:"fetchCheckSec" yaml:"fetchCheckSec"`
	FetchTimeoutSec int64 `json:"fetchTimeoutSec" yaml:"fetchTimeoutSec"`
	RetryCheckSec   int64 `json:"retryCheckSec" yaml:"retryCheckSec"`
	RetryTimeoutSec int64 `json:"retryTimeoutSec" yaml:"retryTimeoutSec"`
}

func (o *Options) RetryCheck() time.Duration {
	return time.Second * time.Duration(o.RetryCheckSec)
}

func (o *Options) FetchCheck() time.Duration {
	return time.Second * time.Duration(o.FetchCheckSec)
}
