package main

import (
	"context"
	"github.com/grpc-boot/base/v2/logger"
	"golang.org/x/exp/rand"
	"time"

	"github.com/grpc-boot/base/v2/gored/mq"
	"github.com/grpc-boot/base/v2/utils"

	"github.com/redis/go-redis/v9"
)

var (
	// 1709629957787-
	// 1709629987799-
	// 1709293904687
	group = `myGroup-new`
	topic = `myTopic`
	red   = redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    "10.16.49.131:6379",
	})
)

func main() {
	err := logger.InitZapWithOption(logger.Option{
		Level: -1,
	})
	if err != nil {
		panic(err)
	}

	go consume(`consumer1`)
	go consume(`consumer2`)
	go consume(`consumer3`)

	go getPending()

	produce()
}

func createMsg() mq.Msg {
	value := map[string]interface{}{
		"id":   time.Now().Unix(),
		"name": time.Now().String(),
	}

	return mq.Msg{
		Topic: topic,
		XMsg: redis.XMessage{
			Values: value,
		},
	}
}

func getPending() {
	consumer, _ := mq.NewConsumer(mq.Option{
		Group:    group,
		Consumer: "info",
	}, red)

	for {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			list, err := consumer.Pending(ctx, topic, time.Second*10, 10)
			if err != nil {
				utils.Red("fetch pending failed with error: %v", err)
				return
			}

			utils.Red("fetch pending list: %+v", list)
			time.Sleep(time.Second * 3)
		}()
	}
}

func produce() {
	producer := mq.NewProducer(mq.Option{
		MaxLength: 8 * 10000,
	}, red)

	for {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			num := 1 + rand.Intn(10)
			msgList := make([]mq.Msg, num)
			for i := 0; i < num; i++ {
				msgList[i] = createMsg()
			}

			if len(msgList) == 1 {
				docId, err := producer.Send(ctx, msgList[0])
				if err != nil {
					utils.Red("produce msg failed with error: %v", err)
				} else {
					utils.Green("produce msg success and got id: %s", docId)
				}
			} else {
				docIdList, err := producer.PipeSend(ctx, msgList...)
				if err != nil {
					utils.Red("produce msg failed with error: %v", err)
				} else {
					utils.Green("produce msg list success and got id list: %v", docIdList)
				}
			}

			time.Sleep(time.Second * 10)
		}()
	}
}

func consume(consumerName string) {
	consumer, _ := mq.NewConsumer(mq.Option{
		Group:    group,
		Consumer: consumerName,
	}, red)

	ch, err := consumer.Consume(topic, 20, time.Second*10, mq.Latest)
	if err != nil {
		utils.RedFatal("consume failed with error: %v", err)
	}

	for {
		msgList, ok := <-ch
		if len(msgList) > 0 {
			for _, msg := range msgList {
				utils.Green("[%s] got msg: %+v", consumerName, msg)
				func(id string) {
					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()

					if rand.Intn(10) == 9 {
						return
					}

					_, err = consumer.Commit(ctx, topic, id)
					if err != nil {
						utils.Red("consumer:%s commit[%s] failed with error: %v", consumerName, id, err)
					}
				}(msg.XMsg.ID)
			}
		}

		if !ok {
			break
		}
	}
}
