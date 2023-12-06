package gored

import (
	"golang.org/x/exp/rand"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	container sync.Map
)

// SetRedis 将*redis.Client放入容器
func SetRedis(key string, options ...redis.Options) {
	list := make([]*redis.Client, len(options))
	for index, opt := range options {
		list[index] = redis.NewClient(&opt)
	}
	container.Store(key, list)
}

// GetRedis 从容器中随机取*redis.Client
func GetRedis(key string) (red *redis.Client, err error) {
	value, exists := container.Load(key)
	if !exists {
		err = ErrNotFound
		return
	}

	list, _ := value.([]*redis.Client)
	if len(list) == 0 {
		err = ErrNotFound
		return
	}

	if len(list) == 1 {
		red = list[0]
		return
	}

	red = list[rand.Intn(len(list))]
	return
}

// GetRedisWithIndex 根据索引从容器中取*redis.Client
func GetRedisWithIndex(key string, index int) (red *redis.Client, err error) {
	value, exists := container.Load(key)
	if !exists {
		err = ErrNotFound
		return
	}

	list, _ := value.([]*redis.Client)
	if len(list) == 0 {
		err = ErrNotFound
		return
	}

	if index >= len(list) || index < 0 {
		err = ErrNotFound
		return
	}

	red = list[index]
	return
}
