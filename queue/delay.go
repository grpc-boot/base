package queue

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/grpc-boot/base/v2/gored"
	"github.com/grpc-boot/base/v2/utils"

	"github.com/redis/go-redis/v9"
)

type Delay interface {
	Set(ctx context.Context, item Item) (err error)
	Done(ctx context.Context, items ...Item) (err error)
	RegisterHandler(handler Handler)
	Start() error
	Stop(timeout time.Duration) error
}

var (
	delayQScript = redis.NewScript(`
		local tKey     = KEYS[1]
		local retryKey = KEYS[2]
		local sMin     = tonumber(ARGV[1])
		local sMax     = tonumber(ARGV[2])
		local tCnt     = tonumber(ARGV[3])
		local tNow     = tonumber(ARGV[4])
		
		local items = redis.call('ZRANGEBYSCORE', tKey, sMin, sMax, 'LIMIT', 0, tCnt)
		if not items[1] then
			return items
		end
		redis.call('ZREM', tKey, unpack(items))

		local rMembers = {}
		for i=1,#items do
			rMembers[2*i -1 ] = tNow
			rMembers[2*i] = items[i]
		end

		redis.call('ZADD', retryKey, unpack(rMembers))
		return items
   `)
)

type delayRedis struct {
	key         string
	retryKey    string
	opt         Options
	handler     Handler
	stop        chan struct{}
	retryTicker *time.Ticker
	fetchTicker *time.Ticker
	red         *redis.Client
}

func NewDelay(key string, red *redis.Client, opt Options) Delay {
	delay := &delayRedis{
		key:         key,
		retryKey:    key + ":retry",
		red:         red,
		opt:         opt,
		stop:        make(chan struct{}),
		fetchTicker: time.NewTicker(opt.FetchCheck()),
		retryTicker: time.NewTicker(opt.RetryCheck()),
	}

	return delay
}

func (dr *delayRedis) autoRetry() {
	for {
		select {
		case <-dr.retryTicker.C:
			go utils.Recover("delay queue auto retry", func(args ...any) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()

				start := time.Now()
				items, _ := dr.retry(ctx, start.Unix(), 100, dr.opt.RetryTimeoutSec)
				fmt.Printf("retry items:%v cost %s\n", items, time.Since(start))
				if len(items) > 0 {
					dr.handler(items)
				}
			})
		case <-dr.stop:
			break
		default:
			time.Sleep(time.Millisecond * 10)
		}
	}
}

func (dr *delayRedis) autoFetch() {
	for {
		select {
		case <-dr.fetchTicker.C:
			go utils.Recover("delay queue auto fetch", func(args ...any) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()

				start := time.Now()
				items, _ := dr.fetch(ctx, start.Unix(), 100, dr.opt.FetchTimeoutSec)
				fmt.Printf("fetch items:%v cost %s\n", items, time.Since(start))
				if len(items) > 0 {
					dr.handler(items)
				}
			})
		case <-dr.stop:
			break
		default:
			time.Sleep(time.Millisecond * 10)
		}
	}
}

func (dr *delayRedis) retry(ctx context.Context, now, count, timeoutSec int64) (items []Item, err error) {
	// lock
	lockCmd := gored.Acquire(ctx, dr.red, dr.retryKey, 8)
	err = gored.DealCmdErr(lockCmd)
	if err != nil || lockCmd.Val() < 1 {
		return
	}

	// lock success
	var (
		token = lockCmd.Val()
		cmd   = dr.red.ZRevRangeByScoreWithScores(ctx, dr.retryKey, &redis.ZRangeBy{
			Min:   strconv.FormatInt(now-timeoutSec, 10),
			Max:   strconv.FormatInt(now-2*timeoutSec, 10),
			Count: count,
		})
	)
	err = gored.DealCmdErr(cmd)
	if err != nil {
		gored.Release(ctx, dr.red, dr.retryKey, token)
		return
	}

	values, _ := cmd.Result()
	if len(values) == 0 {
		gored.Release(ctx, dr.red, dr.retryKey, token)
		return
	}

	// got members
	var (
		members    = make([]redis.Z, 0, len(values))
		delMembers = make([]any, 0, len(values))
	)

	now = time.Now().Unix()
	items = make([]Item, 0, len(values))
	for _, z := range values {
		var item Item
		if utils.JsonDecode(z.Member.(string), &item) != nil {
			continue
		}

		// dead queue
		if item.RetryCount >= dr.opt.MaxRetry {
			members = append(members, redis.Z{
				Score:  0,
				Member: item.Member(),
			})
			continue
		}

		delMembers = append(delMembers, z.Member)

		item.RetryCount++
		members = append(members, redis.Z{
			Score:  float64(now),
			Member: item.Member(),
		})
		items = append(items, item)
	}

	if len(members) == 0 {
		gored.Release(ctx, dr.red, dr.retryKey, token)
		return
	}

	if len(delMembers) == 0 {
		addCmd := dr.red.ZAdd(ctx, dr.retryKey, members...)
		err = gored.DealCmdErr(addCmd)
		if err == nil {
			gored.Release(ctx, dr.red, dr.retryKey, token)
			return
		}

		gored.Release(ctx, dr.red, dr.retryKey, token)
		return nil, err
	}

	pipe := dr.red.TxPipeline()
	pipe.ZRem(ctx, dr.retryKey, delMembers...)
	pipe.ZAdd(ctx, dr.retryKey, members...)
	_, err = pipe.Exec(ctx)
	if err == nil {
		gored.Release(ctx, dr.red, dr.retryKey, token)
		return
	}

	gored.Release(ctx, dr.red, dr.retryKey, token)
	return nil, err
}

func (dr *delayRedis) fetch(ctx context.Context, now, count, timeoutSec int64) (items []Item, err error) {
	var (
		cmd = delayQScript.Run(ctx, dr.red, []string{dr.key, dr.retryKey}, now-timeoutSec, now, count, now)
	)

	err = gored.DealCmdErr(cmd)
	if err != nil {
		return
	}

	value, err := cmd.StringSlice()
	if err != nil || len(value) == 0 {
		return
	}

	items = make([]Item, 0, len(value))
	for _, val := range value {
		var item Item
		if utils.JsonDecode(val, &item) != nil {
			continue
		}
		items = append(items, item)
	}

	return
}

func (dr *delayRedis) RegisterHandler(handler Handler) {
	dr.handler = handler
}

func (dr *delayRedis) Start() error {
	if dr.handler == nil {
		return ErrNoneHandler
	}

	go dr.autoFetch()
	go dr.autoRetry()

	return nil
}

func (dr *delayRedis) Stop(timeout time.Duration) error {
	return utils.Timeout(timeout, func(args ...any) {
		dr.retryTicker.Stop()
		dr.fetchTicker.Stop()
		dr.stop <- struct{}{}
		dr.stop <- struct{}{}
	})
}

func (dr *delayRedis) Done(ctx context.Context, items ...Item) (err error) {
	members := make([]any, len(items))
	for index, item := range items {
		members[index] = item.Member()
	}

	cmd := dr.red.ZRem(ctx, dr.retryKey, members...)

	return gored.DealCmdErr(cmd)
}

func (dr *delayRedis) Set(ctx context.Context, item Item) (err error) {
	cmd := dr.red.ZAdd(ctx, dr.key, redis.Z{
		Score:  float64(item.At),
		Member: item.Member(),
	})

	return gored.DealCmdErr(cmd)
}
