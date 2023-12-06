package gored

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	tokenBucketScript = redis.NewScript(`
		local tKey        = KEYS[1]
		local tCapacity   = tonumber(ARGV[1])
		local current     = tonumber(ARGV[2])
		local rate        = tonumber(ARGV[3])
		local reqNum      = tonumber(ARGV[4])
		local tKeyTimeout = tonumber(ARGV[5])
		
		local bucketInfo = redis.call('HMGET', tKey, 'last_add_time', 'remain_token_num')
		if not bucketInfo[1] then
			redis.call('HMSET', tKey, 'last_add_time', current, 'remain_token_num', tCapacity - 1)
			redis.call('EXPIRE', tKey, tKeyTimeout)
			return 1
		end
		
		local lastAddTime    = tonumber(bucketInfo[1])
		local remainTokenNum = tonumber(bucketInfo[2])
		
		local addTokenNum    = (current - lastAddTime) * rate
		if addTokenNum > 0 then
			lastAddTime    = current
			remainTokenNum = math.min(addTokenNum + remainTokenNum, tCapacity)
			redis.call('HSET', tKey, 'last_add_time', current)
		end
		
		if reqNum > remainTokenNum then
			return 0
		end
		
		redis.call('HSET', tKey, 'remain_token_num', remainTokenNum - reqNum)
		redis.call('EXPIRE', tKey, tKeyTimeout)
		return 1
   `)
)

func GetToken(ctx context.Context, red *redis.Client, key string, current int64, capacity, rate, reqNum, keyTimeoutSecond int) *redis.BoolCmd {
	cmd := tokenBucketScript.Run(ctx, red, []string{key}, capacity, current, rate, reqNum, keyTimeoutSecond)
	val, err := cmd.Int64()

	intCmd := &redis.BoolCmd{}
	intCmd.SetVal(val == 1)
	intCmd.SetErr(err)
	return intCmd
}

func SecondLimitByToken(ctx context.Context, red *redis.Client, key string, limit int, reqNum, burst int) *redis.BoolCmd {
	if burst == 0 {
		burst = limit
	}

	return GetToken(ctx, red, key, time.Now().Unix(), burst, limit, reqNum, 5)
}
