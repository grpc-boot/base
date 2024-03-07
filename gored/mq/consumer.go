package mq

import (
	"context"
	"sort"
	"time"

	"github.com/grpc-boot/base/v2/logger"
	"github.com/grpc-boot/base/v2/utils"

	"github.com/redis/go-redis/v9"
)

const (
	Earliest = `0`
	Latest   = `$`
)

const (
	peddingId     = `0`
	noDeliveredId = ">"
)

func (mq *Mq) Group(ctx context.Context, startId string) (err error) {
	statusCmd := mq.pool.XGroupCreateMkStream(ctx, mq.option.ConsumerTopic, mq.option.Group, startId)
	return statusCmd.Err()
}

func (mq *Mq) Close(timeout time.Duration) error {
	if mq.consumeChan == nil {
		return nil
	}

	if !mq.done.CompareAndSwap(false, true) {
		return nil
	}

	done := make(chan struct{}, 1)
	go func() {
		for !mq.balanceDone.Load() || !mq.readDone.Load() || len(mq.consumeChan) > 0 {
			time.Sleep(time.Millisecond * 100)
		}
		done <- struct{}{}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}

func (mq *Mq) Consume(maxCount int64, blockTime time.Duration, startId string) (msgChan <-chan []Msg, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err = mq.Group(ctx, startId)
	if IsErrBusyGroup(err) {
		err = nil
	}

	if err != nil {
		return
	}

	mq.consumeChan = make(chan []Msg, mq.option.ChanSize)

	go func() {
		id := peddingId
		for {
			if mq.done.Load() {
				mq.readDone.Store(true)
				return
			}

			cmd := mq.pool.XReadGroup(context.Background(), &redis.XReadGroupArgs{
				Group:    mq.option.Group,
				Consumer: mq.option.Consumer,
				Streams:  []string{mq.option.ConsumerTopic, id},
				Count:    maxCount,
				Block:    blockTime,
				NoAck:    mq.option.AutoCommit,
			})

			streamData, er := cmd.Result()
			if er != nil {
				logger.ZapError("consume failed",
					logger.Error(er),
					logger.Topic(mq.option.ConsumerTopic),
				)

				time.Sleep(mq.option.RetryDuration())
				continue
			}

			if len(streamData) == 0 {
				continue
			}

			for _, msgData := range streamData {
				if len(msgData.Messages) == 0 {
					if id == peddingId {
						id = noDeliveredId
					}
					continue
				}

				msgList := make([]Msg, len(msgData.Messages))
				for index, msg := range msgData.Messages {
					msgList[index] = Msg{
						Topic: msgData.Stream,
						XMsg:  msg,
					}
				}

				mq.consumeChan <- msgList
			}
		}
	}()

	go mq.autoBalance()

	return mq.consumeChan, nil
}

func (mq *Mq) Commit(ctx context.Context, topic string, idList ...string) (okCount int64, err error) {
	cmd := mq.pool.XAck(ctx, topic, mq.option.Group, idList...)
	return cmd.Result()
}

func (mq *Mq) autoBalance() {
	for {
		if mq.done.Load() {
			mq.balanceDone.Store(true)
			return
		}

		time.Sleep(time.Second * 10)

		if mq.done.Load() {
			mq.balanceDone.Store(true)
			return
		}

		err := mq.balancer(context.Background())
		if err != nil {
			logger.ZapError("auto balance failed",
				logger.Error(err),
				logger.Topic(mq.option.ConsumerTopic),
			)
		}
	}
}

func (mq *Mq) balancer(ctx context.Context) (err error) {
	list, err := mq.pool.XInfoConsumers(ctx, mq.option.ConsumerTopic, mq.option.Group).Result()
	if err != nil {
		return err
	}

	if len(list) == 0 {
		return nil
	}

	var (
		canBalanceList  = make([]redis.XInfoConsumer, 0)
		needBalanceList = make([]redis.XInfoConsumer, 0)
	)

	for _, infoConsumer := range list {
		if infoConsumer.Pending == 0 {
			canBalanceList = append(canBalanceList, infoConsumer)
			continue
		}
		needBalanceList = append(needBalanceList, infoConsumer)
	}

	if len(needBalanceList) == 0 {
		return
	}

	if len(canBalanceList) == 0 {
		canBalanceList = list
	}

	sort.SliceStable(canBalanceList, func(i, j int) bool {
		if list[i].Pending == list[j].Pending {
			return list[i].Inactive < list[j].Inactive
		}
		return list[i].Pending < list[j].Pending
	})

	var (
		msgList []redis.XMessage
		i             = 1
		size    int64 = 128
		start         = `0-0`
	)

	for _, infoConsumer := range needBalanceList {
		num := utils.Ceil[int](float64(infoConsumer.Pending) / float64(size))
		for n := 0; n < num; n++ {
			i++
			if start == "" {
				start = `0-0`
			}

			cmd := mq.pool.XAutoClaim(ctx, &redis.XAutoClaimArgs{
				Stream:   mq.option.ConsumerTopic,
				Group:    mq.option.Group,
				MinIdle:  mq.option.MsgMinIdle(),
				Start:    start,
				Count:    size,
				Consumer: canBalanceList[i%len(canBalanceList)].Name,
			})

			msgList, start, err = cmd.Result()
			if err != nil {
				logger.ZapError("autoclaim msg failed",
					logger.Error(err),
					logger.Cmd(cmd.String()),
				)
			} else {
				logger.ZapInfo("autoclaim msg",
					logger.Cmd(cmd.String()),
				)

				if len(msgList) > 0 {
					ml := make([]Msg, len(msgList))
					for index, m := range msgList {
						ml[index] = Msg{Topic: mq.option.ConsumerTopic, XMsg: m}
					}
					mq.consumeChan <- ml
				}
			}
		}
	}

	return nil
}
