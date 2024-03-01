package mq

import (
	"context"
	"time"

	"github.com/grpc-boot/base/v2/logger"

	"github.com/redis/go-redis/v9"
)

const (
	peddingId     = `0`
	noDeliveredId = ">"
)

func (mq *Mq) Consume(topic string, maxCount int64, blockTime time.Duration) (data <-chan []Msg, err error) {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
		statusCmd   = mq.pool.XGroupCreateMkStream(ctx, topic, mq.option.Group, "0")
	)

	defer cancel()

	err = statusCmd.Err()
	if IsErrBusyGroup(err) {
		err = nil
	}

	if err != nil {
		return
	}

	dataCh := make(chan []Msg, 1024)

	go func() {
		id := peddingId
		for {
			cmd := mq.pool.XReadGroup(context.Background(), &redis.XReadGroupArgs{
				Group:    mq.option.Group,
				Consumer: mq.option.Consumer,
				Streams:  []string{topic, id},
				Count:    maxCount,
				Block:    blockTime,
				NoAck:    mq.option.AutoCommit,
			})

			streamData, er := cmd.Result()
			if er != nil {
				logger.ZapError("consume failed",
					logger.Error(er),
					logger.Topic(topic),
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

				dataCh <- msgList
			}
		}
	}()

	return dataCh, nil
}

func (mq *Mq) Commit(ctx context.Context, topic string, idList ...string) (okCount int64, err error) {
	cmd := mq.pool.XAck(ctx, topic, mq.option.Group, idList...)
	return cmd.Result()
}

func (mq *Mq) compensate(ctx context.Context, topic string, minIdleDuration time.Duration) (err error) {
	for {
		list, er := mq.PendingWithConsumer(ctx, topic, minIdleDuration, 100)
		if er != nil || len(list) == 0 {
			return er
		}

		idList := make([]string, 0, len(list))
		for _, pending := range list {
			if pending.RetryCount >= mq.option.MaxRetryCount {
				continue
			}
			idList = append(idList, pending.ID)
		}

		mq.pool.XClaim(ctx, &redis.XClaimArgs{
			Stream:   topic,
			Group:    mq.option.Group,
			Consumer: mq.option.Consumer,
			MinIdle:  minIdleDuration,
			Messages: idList,
		})
	}
}

func (mq *Mq) Balance(ctx context.Context, topic string, minIdleDuration time.Duration) (err error) {
	info, err := mq.FullInfo(ctx, topic, 1)
	if err != nil {
		return err
	}

	if len(info.Groups) == 0 {
		return nil
	}

	for _, group := range info.Groups {
		if group.Name == mq.option.Group {
			length := len(group.Consumers)
			if length == 0 {
				return nil
			}
		}
	}

	return nil
}
