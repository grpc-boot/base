package gored

import (
	"time"

	"github.com/redis/go-redis/v9"
)

var DefaultOptions = func() redis.Options {
	return redis.Options{
		DB:                    0,
		MaxRetries:            3,
		DialTimeout:           time.Second,
		ReadTimeout:           time.Second * 3,
		WriteTimeout:          time.Second * 3,
		ContextTimeoutEnabled: true,
		PoolSize:              16,
		MinIdleConns:          1,
		MaxIdleConns:          2,
		MaxActiveConns:        10,
		ConnMaxIdleTime:       time.Second * 60,
		ConnMaxLifetime:       time.Second * 300,
	}
}
