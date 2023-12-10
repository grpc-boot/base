package gored

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/grpc-boot/base/v2/kind"
	"github.com/grpc-boot/base/v2/utils"

	"github.com/redis/go-redis/v9"
	"go.uber.org/atomic"
)

const (
	MinInterval = time.Second * 3
)

type Conf struct {
	key      string
	red      *redis.Client
	interval time.Duration
	value    atomic.Value
}

func NewConf(key string, interval time.Duration, red *redis.Client) *Conf {
	config := &Conf{
		key:      key,
		interval: interval,
		red:      red,
	}

	config.sync()

	go config.autoSync()

	return config
}

func SetConf[T kind.RedisValue](c *Conf, key string, value T) (isNew bool, err error) {
	var cmd *redis.IntCmd
	TimeoutDo(time.Second*3, func(ctx context.Context) {
		cmd = c.red.HSet(ctx, c.key, key, value)
	})

	err = DealCmdErr(cmd)
	if err == nil {
		isNew = cmd.Val() == 1
	}

	return
}

func (c *Conf) String(key, defaultValue string) string {
	m, _ := c.value.Load().(map[string]string)
	value, exists := m[key]
	if !exists {
		return defaultValue
	}

	return value
}

func (c *Conf) Int(key string, defaultValue int64) int64 {
	m, _ := c.value.Load().(map[string]string)
	val, exists := m[key]
	if !exists {
		return defaultValue
	}

	value, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return defaultValue
	}

	return value
}

func (c *Conf) Uint(key string, defaultValue uint64) uint64 {
	m, _ := c.value.Load().(map[string]string)
	val, exists := m[key]
	if !exists {
		return defaultValue
	}

	value, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return defaultValue
	}

	return value
}

func (c *Conf) Float(key string, defaultValue float64) float64 {
	m, _ := c.value.Load().(map[string]string)
	val, exists := m[key]
	if !exists {
		return defaultValue
	}

	value, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return defaultValue
	}

	return value
}

func (c *Conf) Bool(key string, defaultValue bool) bool {
	m, _ := c.value.Load().(map[string]string)
	val, exists := m[key]
	if !exists {
		return defaultValue
	}

	value, err := strconv.ParseBool(val)
	if err != nil {
		return defaultValue
	}

	return value
}

func (c *Conf) autoSync() {
	if c.interval < MinInterval {
		c.interval = MinInterval
	}

	ticker := time.NewTicker(c.interval)
	for range ticker.C {
		utils.Recover(fmt.Sprintf("sync redis conf[%s]", c.key), func(args ...any) {
			c.sync()
		})
	}
}

func (c *Conf) sync() {
	var cmd *redis.MapStringStringCmd
	TimeoutDo(time.Second*3, func(ctx context.Context) {
		cmd = c.red.HGetAll(ctx, c.key)
	})

	if err := DealCmdErr(cmd); err != nil {
		return
	}

	c.value.Store(cmd.Val())
}
