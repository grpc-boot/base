package mq

import "github.com/redis/go-redis/v9"

type Data map[string]string

type Msg struct {
	Topic string
	XMsg  redis.XMessage
}
