package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-boot/base/v2/gored"
	"github.com/grpc-boot/base/v2/gored/mq"
	"github.com/grpc-boot/base/v2/logger"
	"github.com/grpc-boot/base/v2/utils"

	"github.com/redis/go-redis/v9"
)

var (
	red *redis.Client
)

func main() {
	err := logger.InitZapWithOption(logger.Option{
		Level: -1,
	})
	if err != nil {
		panic(err)
	}

	redOptions := gored.DefaultOptions()
	redOptions.Addr = "10.16.49.131:6379"
	red = redis.NewClient(&redOptions)

	consumer1, _ := mq.NewConsumer(mq.Option{
		Group:         `myGroup-new`,
		ConsumerTopic: `base-topic`,
	}, red)

	ch1, err := consumer1.Consume(20, time.Second*20, mq.Earliest)
	// 最新
	//ch1, err := consumer1.Consume(20, time.Second*20, mq.Latest)
	if err != nil {
		utils.RedFatal("consume msg failed with error: %s", err)
	}

	go consume(consumer1, ch1)

	var sig = make(chan os.Signal, 1)
	signal.Notify(sig)

	for {
		val := <-sig
		switch val {
		case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM:
			err = consumer1.Close(time.Second * 10)
			if err != nil {
				logger.ZapError("close consumer failed",
					logger.Error(err),
					logger.Consumer(`consumer1`),
				)
			}
			return
		default:
			continue
		}
	}
}

func consume(consumer *mq.Mq, ch <-chan []mq.Msg) {
	for msgList := range ch {
		if len(msgList) > 0 {
			for _, msg := range msgList {
				utils.Green("got msg: %+v", msg)
				func(m redis.XMessage) {
					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()

					event := mq.Msg2Event(m)
					if event.Id != "" && event.Name != "" {
						utils.Green("got event: %+v", event)
					}

					_, err := consumer.Commit(ctx, msg.Topic, m.ID)
					if err != nil {
						utils.Red("commit[%s] failed with error: %v", m.ID, err)
					}
				}(msg.XMsg)
			}
		}
	}
}
