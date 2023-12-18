package http_client

import "time"

var (
	DefaultOptions = func() Options {
		return Options{
			MaxIdleConns:        8,
			MaxIdleConnsPerHost: 2,
			MaxConnsPerHost:     8,
			TimeoutSec:          30,
			IdleConnTimeoutSec:  60,
			WriteBufferSize:     4 << 10,
			ReadBufferSize:      4 << 10,
		}
	}
)

type Options struct {
	DisableKeepAlives   bool `json:"disableKeepAlives" yaml:"disableKeepAlives"`
	DisableCompression  bool `json:"disableCompression" yaml:"disableCompression"`
	MaxIdleConns        int  `json:"maxIdleConns" yaml:"maxIdleConns"`
	MaxIdleConnsPerHost int  `json:"maxIdleConnsPerHost" yaml:"maxIdleConnsPerHost"`
	MaxConnsPerHost     int  `json:"maxConnsPerHost" yaml:"maxConnsPerHost"`
	TimeoutSec          int  `json:"timeoutSec" yaml:"timeoutSec"`
	IdleConnTimeoutSec  int  `json:"idleConnTimeoutSec" yaml:"idleConnTimeoutSec"`
	WriteBufferSize     int  `json:"writeBufferSize" yaml:"writeBufferSize"`
	ReadBufferSize      int  `json:"readBufferSize" yaml:"readBufferSize"`
}

func (o *Options) Timeout() time.Duration {
	return time.Second * time.Duration(o.TimeoutSec)
}

func (o *Options) IdleConnTimeout() time.Duration {
	return time.Second * time.Duration(o.IdleConnTimeoutSec)
}
