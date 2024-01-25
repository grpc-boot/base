package queue

import "time"

var (
	DefaultOptions = func() Options {
		return Options{
			MaxRetry:        16,
			FetchCheckSec:   5,
			FetchForwardSec: 60,
			FetchPerCount:   128,
			RetryCheckSec:   60,
			RetryForwardSec: 600,
			RetryPerCount:   1024,
		}
	}
)

type Options struct {
	MaxRetry        int64 `json:"maxRetry" yaml:"maxRetry"`
	FetchCheckSec   int64 `json:"fetchCheckSec" yaml:"fetchCheckSec"`
	FetchPerCount   int64 `json:"fetchPerCount" yaml:"fetchPerCount"`
	FetchForwardSec int64 `json:"fetchForwardSec" yaml:"fetchForwardSec"`
	RetryCheckSec   int64 `json:"retryCheckSec" yaml:"retryCheckSec"`
	RetryPerCount   int64 `json:"retryPerCount" yaml:"retryPerCount"`
	RetryForwardSec int64 `json:"retryForwardSec" yaml:"retryForwardSec"`
}

func (o *Options) RetryForwardEndAt(now int64) int64 {
	if o.RetryForwardSec < 1 {
		return now
	}

	return now - o.RetryForwardSec
}

func (o *Options) FetchForwardBeginAt(now int64) int64 {
	if o.FetchForwardSec < 1 {
		return 1
	}

	return now - o.FetchForwardSec
}

func (o *Options) RetryCheck() time.Duration {
	return time.Second * time.Duration(o.RetryCheckSec)
}

func (o *Options) FetchCheck() time.Duration {
	return time.Second * time.Duration(o.FetchCheckSec)
}
