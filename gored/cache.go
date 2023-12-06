package gored

import (
	"context"
	"fmt"
	"time"

	"github.com/grpc-boot/base/v2/gored/cache"
	"github.com/grpc-boot/base/v2/kind/msg"
	"github.com/grpc-boot/base/v2/logger"

	"github.com/redis/go-redis/v9"
)

const (
	cacheKeyFormat = `gored_C:%s`
	lockSeconds    = 8
	cacheTimeout   = time.Duration(3600*24*7) * time.Second
)

func GetItemWithCacheTimeout[V msg.Value](timeout time.Duration, red *redis.Client, key string, current, cacheSeconds int64, handler msg.Handler[V]) (item cache.Item, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return GetItemWithCache(ctx, red, key, current, cacheSeconds, handler)
}

func GetItemWithCache[V msg.Value](ctx context.Context, red *redis.Client, key string, current, cacheSeconds int64, handler msg.Handler[V]) (item cache.Item, err error) {
	key = fmt.Sprintf(cacheKeyFormat, key)
	cmd := red.Get(ctx, key)
	err = DealCmdErr(cmd)
	if err != nil {
		return
	}

	redisValue, _ := cmd.Bytes()

	//redis中没有数据
	if len(redisValue) == 0 {
		err = updateCache(ctx, red, key, &item, current, handler)
		return item, err
	}

	var left []byte
	left, err = item.UnmarshalMsg(redisValue)
	if err != nil || len(left) > 0 {
		err = updateCache(ctx, red, key, &item, current, handler)
		return
	}

	if item.UpdatedAt == 0 {
		err = updateCache(ctx, red, key, &item, current, handler)
		return item, err
	}

	//缓存有效
	if item.Hit(cacheSeconds, current) {
		return item, err
	}

	//-------------------缓存失效-----------------------
	//去拿锁
	acqCmd := Acquire(ctx, red, key, lockSeconds)
	err = DealCmdErr(acqCmd)
	//未获得锁
	if err != nil || acqCmd.Val() == 0 {
		return item, nil
	}

	// 获得锁
	err = updateCache(ctx, red, key, &item, current, handler)
	if err == nil {
		_ = Release(ctx, red, key, acqCmd.Val())
	}

	return item, err
}

func updateCache[V msg.Value](ctx context.Context, red *redis.Client, key string, item *cache.Item, current int64, handler msg.Handler[V]) (err error) {
	item.Value, err = handler()
	if err != nil {
		logger.ZapError("cache exec handler failed",
			logger.Key(key),
			logger.Error(err),
		)
		return
	}

	item.PackValue()

	item.UpdatedAt = current
	if item.CreatedAt < 1 {
		item.CreatedAt = current
	}

	data, err := item.Marshal()
	if err != nil {
		err = ErrInvalidDataType
		return
	}

	cmd := red.SetEx(ctx, key, data, cacheTimeout)
	err = DealCmdErr(cmd)
	return
}
