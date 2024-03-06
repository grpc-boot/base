package main

import (
	"context"
	"golang.org/x/exp/rand"
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
	group = `myGroup-new`
	topic = `myTopic`
	red   *redis.Client
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
		Group:         group,
		Consumer:      `consumer1`,
		ConsumerTopic: topic,
	}, red)

	ch1, err := consumer1.Consume(20, time.Second*20, mq.Earliest)
	//ch1, err := consumer1.Consume(20, time.Second*20, mq.Latest)
	if err != nil {
		utils.RedFatal("consume msg failed with error: %s", err)
	}

	go consume(consumer1, ch1, `consumer1`)

	go getPending()

	go produce()

	var sig = make(chan os.Signal, 1)
	signal.Notify(sig)

	for {
		val := <-sig
		switch val {
		case syscall.SIGINT:
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
		Group:         group,
		Consumer:      "info",
		ConsumerTopic: topic,
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

func consume(consumer *mq.Mq, ch <-chan []mq.Msg, consumerName string) {
	for msgList := range ch {
		if len(msgList) > 0 {
			for _, msg := range msgList {
				utils.Green("[%s] got msg: %+v", consumerName, msg)
				func(id string) {
					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()

					if rand.Intn(10) == 9 {
						return
					}

					_, err := consumer.Commit(ctx, topic, id)
					if err != nil {
						utils.Red("consumer:%s commit[%s] failed with error: %v", consumerName, id, err)
					}
				}(msg.XMsg.ID)
			}
		}
	}

	logger.ZapError("chan has closed")
}
