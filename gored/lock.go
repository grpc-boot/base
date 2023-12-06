package gored

import (
	"context"
	"fmt"
	"golang.org/x/exp/rand"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	lockFormat = `gored_L:%s`
)

var (
	delTokenScript = redis.NewScript(`if redis.call('get', KEYS[1]) == ARGV[1]
	then
		return redis.call('del', KEYS[1])
	end
	return 0`)
)

func init() {
	rand.Seed(uint64(time.Now().UnixNano()))
}

func AcquireWithTimeout(timeout time.Duration, red *redis.Client, key string, lockSeconds int64) (intCmd *redis.IntCmd) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return Acquire(ctx, red, key, lockSeconds)
}

// Acquire 加锁
func Acquire(ctx context.Context, red *redis.Client, key string, lockSeconds int64) (intCmd *redis.IntCmd) {
	intCmd = &redis.IntCmd{}

	t := rand.Int63()

	var (
		cmd = red.SetNX(ctx, fmt.Sprintf(lockFormat, key), t, time.Duration(lockSeconds)*time.Second)
	)

	if cmd.Err() != nil {
		intCmd.SetErr(cmd.Err())
		return
	}

	if cmd.Val() {
		intCmd.SetVal(t)
	}
	return
}

func ReleaseWithTimeout(timeout time.Duration, red *redis.Client, key string, token int64) (intCmd *redis.IntCmd) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return Release(ctx, red, key, token)
}

// Release 释放锁
func Release(ctx context.Context, red *redis.Client, key string, token int64) (intCmd *redis.IntCmd) {
	intCmd = &redis.IntCmd{}

	var (
		cmd = delTokenScript.Run(ctx, red, []string{fmt.Sprintf(lockFormat, key)}, token)
	)

	if cmd.Err() != nil {
		intCmd.SetErr(cmd.Err())
		return
	}

	delNum, err := cmd.Int64()
	intCmd.SetErr(err)
	intCmd.SetVal(delNum)

	return
}
