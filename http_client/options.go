package http_client

import "time"

var (
	DefaultOptions = func() Options {
		return Options{
			MaxIdleConns:        8,
			MaxIdleConnsPerHost: 2,
			MaxConnsPerHost:     8,
			TimeoutSec:          30,
			DialTimeoutSec:      30,
			IdleConnTimeoutSec:  60,
			KeepaliveSec:        60,
			WriteBufferSize:     4 << 10,
			ReadBufferSize:      4 << 10,
		}
	}
)

type Options struct {
	MaxIdleConns        int  `json:"maxIdleConns" yaml:"maxIdleConns"`
	MaxIdleConnsPerHost int  `json:"maxIdleConnsPerHost" yaml:"maxIdleConnsPerHost"`
	MaxConnsPerHost     int  `json:"maxConnsPerHost" yaml:"maxConnsPerHost"`
	TimeoutSec          int  `json:"timeoutSec" yaml:"timeoutSec"`
	DialTimeoutSec      int  `json:"dialTimeoutSec" yaml:"dialTimeoutSec"`
	IdleConnTimeoutSec  int  `json:"idleConnTimeoutSec" yaml:"idleConnTimeoutSec"`
	DisableKeepAlives   bool `json:"disableKeepAlives" yaml:"disableKeepAlives"`
	DisableCompression  bool `json:"disableCompression" yaml:"disableCompression"`
	KeepaliveSec        int  `json:"keepaliveSec" yaml:"keepaliveSec"`
	WriteBufferSize     int  `json:"writeBufferSize" yaml:"writeBufferSize"`
	ReadBufferSize      int  `json:"readBufferSize" yaml:"readBufferSize"`
}

func (o *Options) Timeout() time.Duration {
	return time.Second * time.Duration(o.TimeoutSec)
}

func (o *Options) DialTimeout() time.Duration {
	return time.Duration(o.DialTimeoutSec) * time.Second
}

func (o *Options) KeepaliveTime() time.Duration {
	return time.Duration(o.KeepaliveSec) * time.Second
}

func (o *Options) IdleConnTimeout() time.Duration {
	return time.Second * time.Duration(o.IdleConnTimeoutSec)
}
