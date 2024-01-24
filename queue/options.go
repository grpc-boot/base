package queue

import "time"

var (
	DefaultOptions = func() Options {
		return Options{
			MaxRetry:        16,
			FetchCheckSec:   5,
			FetchForwardSec: 60,
			RetryCheckSec:   60,
			RetryForwardSec: 600,
		}
	}
)

type Options struct {
	MaxRetry        int64 `json:"maxRetry" yaml:"maxRetry"`
	FetchCheckSec   int64 `json:"fetchCheckSec" yaml:"fetchCheckSec"`
	FetchForwardSec int64 `json:"fetchForwardSec" yaml:"fetchForwardSec"`
	RetryCheckSec   int64 `json:"retryCheckSec" yaml:"retryCheckSec"`
	RetryForwardSec int64 `json:"retryForwardSec" yaml:"retryForwardSec"`
}

func (o *Options) RetryCheck() time.Duration {
	return time.Second * time.Duration(o.RetryCheckSec)
}

func (o *Options) FetchCheck() time.Duration {
	return time.Second * time.Duration(o.FetchCheckSec)
}
