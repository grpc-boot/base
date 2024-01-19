package elasticsearch

import (
	"time"

	"github.com/grpc-boot/base/v2/http_client"
	"github.com/grpc-boot/base/v2/utils"
)

func DefaultOption() Options {
	return Options{
		MaxOpenConns:       8,
		MaxIdleConns:       2,
		ConnMaxLifetimeSec: 600,
		ConnMaxIdleTimeSec: 60,
	}
}

func OptionsWithJson(conf string) (opt Options, err error) {
	opt = DefaultOption()
	err = utils.JsonDecode(conf, &opt)
	return
}

func OptionsWithYaml(conf string) (opt Options, err error) {
	opt = DefaultOption()
	err = utils.YamlDecode(conf, &opt)
	return
}

type Options struct {
	BaseUrl            string `json:"baseUrl" yaml:"baseUrl"`
	UserName           string `json:"userName" yaml:"userName"`
	Password           string `json:"password" yaml:"password"`
	MaxIdleConns       int    `json:"maxIdleConns" yaml:"maxIdleConns"`
	MaxOpenConns       int    `json:"maxOpenConns" yaml:"maxOpenConns"`
	TimeoutSec         int    `json:"timeoutSec" yaml:"timeoutSec"`
	DialTimeoutSec     int    `json:"dialTimeoutSec" yaml:"dialTimeoutSec"`
	KeepaliveSec       int    `json:"keepaliveSec" yaml:"keepaliveSec"`
	ConnMaxIdleTimeSec int    `json:"connMaxIdleTimeSec" yaml:"connMaxIdleTimeSec"`
	ConnMaxLifetimeSec int    `json:"connMaxLifetimeSec" yaml:"connMaxLifetimeSec"`
	ReadBufferSize     int    `json:"readBufferSize" yaml:"readBufferSize"`
	WriteBufferSize    int    `json:"writeBufferSize" yaml:"writeBufferSize"`
}

func (o *Options) httpOptions() http_client.Options {
	opt := http_client.DefaultOptions()
	if o.MaxIdleConns > 0 {
		opt.MaxIdleConns = o.MaxIdleConns
		opt.MaxConnsPerHost = o.MaxIdleConns
	}

	if o.MaxOpenConns > 0 {
		opt.MaxConnsPerHost = o.MaxOpenConns
	}

	if o.TimeoutSec > 0 {
		opt.TimeoutSec = o.TimeoutSec
	}

	if o.DialTimeoutSec > 0 {
		opt.DialTimeoutSec = o.DialTimeoutSec
	}

	if o.ConnMaxIdleTimeSec > 0 {
		opt.IdleConnTimeoutSec = o.ConnMaxIdleTimeSec
	}

	if o.KeepaliveSec > 0 {
		opt.KeepaliveSec = o.KeepaliveSec
	}

	if o.WriteBufferSize > 0 {
		opt.WriteBufferSize = o.WriteBufferSize
	}

	if o.ReadBufferSize > 0 {
		opt.ReadBufferSize = o.ReadBufferSize
	}

	return opt
}

func (o *Options) Timeout() time.Duration {
	return time.Duration(o.TimeoutSec) * time.Second
}

func (o *Options) DialTimeout() time.Duration {
	return time.Duration(o.DialTimeoutSec) * time.Second
}

func (o *Options) KeepaliveTime() time.Duration {
	return time.Duration(o.KeepaliveSec) * time.Second
}

func (o *Options) ConnMaxIdleTime() time.Duration {
	return time.Duration(o.ConnMaxIdleTimeSec) * time.Second
}

func (o *Options) ConnMaxLifetime() time.Duration {
	return time.Duration(o.ConnMaxLifetimeSec) * time.Second
}
