package mq

import (
	"context"

	"github.com/redis/go-redis/v9"
)

const (
	autoCreateId = `*`
)

func (mq *Mq) Send(ctx context.Context, msg Msg) (docId string, err error) {
	id := autoCreateId
	if msg.XMsg.ID != "" {
		id = msg.XMsg.ID
	}

	cmd := mq.pool.XAdd(ctx, &redis.XAddArgs{
		Stream:     msg.Topic,
		NoMkStream: true,
		MaxLen:     mq.option.MaxLength,
		ID:         id,
		Values:     msg.XMsg.Values,
		Approx:     true,
	})

	return cmd.Result()
}

func (mq *Mq) PipeSend(ctx context.Context, msgList ...Msg) (docIdList []string, err error) {
	if len(msgList) == 1 {
		docId, er := mq.Send(ctx, msgList[0])
		if er != nil {
			return nil, er
		}
		return []string{docId}, nil
	}

	pipe := mq.pool.Pipeline()

	for _, msg := range msgList {
		id := autoCreateId
		if msg.XMsg.ID != "" {
			id = msg.XMsg.ID
		}

		pipe.XAdd(ctx, &redis.XAddArgs{
			Stream:     msg.Topic,
			NoMkStream: true,
			MaxLen:     mq.option.MaxLength,
			ID:         id,
			Values:     msg.XMsg.Values,
			Approx:     true,
		})
	}

	cmdList, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	docIdList = make([]string, len(cmdList))
	for index, cmd := range cmdList {
		docIdList[index], _ = cmd.(*redis.StringCmd).Result()
	}

	return
}
