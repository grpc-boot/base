package monitor

import (
	"context"
	"time"

	"github.com/grpc-boot/base/v2/gored"

	"github.com/redis/go-redis/v9"
)

var (
	PipeTimeout     = time.Second * 3
	StorageLifetime = time.Hour * 24 * 7
)

type Storage func(info *MonitorInfo) error

func RedisStorage(red *redis.Client) Storage {
	return func(info *MonitorInfo) (err error) {
		var (
			now     = time.Now()
			dateStr = now.Format("20060102")
			timeStr = now.Format("15:04")
		)

		if len(info.GaugesInfo) > 0 {
			gored.TimeoutDo(PipeTimeout, func(ctx context.Context) {
				pipe := red.Pipeline()
				for _, gauge := range info.GaugesInfo {
					key := gauge.Key(info.Name, dateStr)

					pipe.HIncrBy(
						ctx,
						key,
						timeStr,
						int64(gauge.Value),
					)

					if timeStr[3:] == "00" {
						pipe.Expire(ctx, key, StorageLifetime)
					}
				}

				_, err = pipe.Exec(ctx)
				cmd := &redis.BoolCmd{}
				cmd.SetErr(err)
				cmd.SetVal(err == nil)
				_ = gored.DealCmdErr(cmd)
			})
		}

		if len(info.CodesInfo) > 0 {
			for _, group := range info.CodesInfo {
				if len(group) < 1 {
					continue
				}

				for _, gauge := range group {
					gored.TimeoutDo(PipeTimeout, func(ctx context.Context) {
						var (
							pipe = red.Pipeline()
							key  = gauge.Key(info.Name, dateStr)
						)

						pipe.HIncrBy(
							ctx,
							key,
							timeStr,
							int64(gauge.Value),
						)

						if len(gauge.Sub) > 0 {
							for _, sub := range gauge.Sub {
								pipe.HIncrBy(
									ctx,
									key,
									sub.Field(timeStr),
									int64(sub.Value),
								)
							}
						}

						if timeStr[3:] == "00" {
							pipe.Expire(ctx, key, StorageLifetime)
						}

						_, err = pipe.Exec(ctx)
						cmd := &redis.BoolCmd{}
						cmd.SetErr(err)
						cmd.SetVal(err == nil)
						_ = gored.DealCmdErr(cmd)
					})
				}
			}
		}

		return
	}
}
