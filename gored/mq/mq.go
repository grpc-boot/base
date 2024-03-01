package mq

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Mq struct {
	option Option
	pool   *redis.Client
}

func NewConsumer(opt Option, pool *redis.Client) (*Mq, error) {
	option := formatOption(opt)

	if option.Group == "" {
		return nil, ErrGroupEmpty
	}

	if option.Consumer == "" {
		return nil, ErrConsumerEmpty
	}

	return &Mq{
		option: option,
		pool:   pool,
	}, nil
}

func NewProducer(opt Option, pool *redis.Client) *Mq {
	return &Mq{
		option: opt,
		pool:   pool,
	}
}

func (mq *Mq) Info(ctx context.Context, topic string) (info *redis.XInfoStream, err error) {
	cmd := mq.pool.XInfoStream(ctx, topic)
	return cmd.Result()
}

func (mq *Mq) FullInfo(ctx context.Context, topic string, count int) (info *redis.XInfoStreamFull, err error) {
	cmd := mq.pool.XInfoStreamFull(ctx, topic, count)
	return cmd.Result()
}

func (mq *Mq) Pending(ctx context.Context, topic string, minIdleDuration time.Duration, count int64) (list []redis.XPendingExt, err error) {
	cmd := mq.pool.XPendingExt(ctx, &redis.XPendingExtArgs{
		Stream: topic,
		Group:  mq.option.Group,
		Idle:   minIdleDuration,
		Start:  "-",
		End:    "+",
		Count:  count,
	})

	return cmd.Result()
}

func (mq *Mq) PendingWithConsumer(ctx context.Context, topic string, minIdleDuration time.Duration, count int64) (list []redis.XPendingExt, err error) {
	cmd := mq.pool.XPendingExt(ctx, &redis.XPendingExtArgs{
		Stream:   topic,
		Group:    mq.option.Group,
		Idle:     minIdleDuration,
		Start:    "-",
		End:      "+",
		Count:    count,
		Consumer: mq.option.Consumer,
	})

	return cmd.Result()
}
