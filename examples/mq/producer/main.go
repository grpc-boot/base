package main

import (
	"context"
	"golang.org/x/exp/rand"
	"strconv"
	"time"

	"github.com/grpc-boot/base/v2/gored"
	"github.com/grpc-boot/base/v2/gored/mq"
	"github.com/grpc-boot/base/v2/logger"
	"github.com/grpc-boot/base/v2/utils"

	"github.com/redis/go-redis/v9"
)

var (
	eventList = []string{"login", "register", "login-out", "pay"}
	red       *redis.Client
)

func main() {
	err := logger.InitZapWithOption(logger.Option{
		Level: -1,
	})
	if err != nil {
		panic(err)
	}

	redOptions := gored.DefaultOptions()
	redOptions.Addr = "127.0.0.1:6379"
	red = redis.NewClient(&redOptions)

	producer := mq.NewProducer(mq.Option{
		MaxLength: 8 * 10000,
	}, red)

	produce(producer, `base-topic`)
}

func createEvent() mq.Event {
	val, _ := utils.JsonEncode(map[string]any{
		"createdAt":       time.Now().Unix(),
		"createdDateTime": time.Now().String(),
	})

	return mq.Event{
		Name: eventList[rand.Intn(len(eventList))],
		Id:   strconv.FormatInt(time.Now().UnixNano(), 10),
		Data: val,
	}
}

func produce(producer *mq.Mq, topic string) {
	for {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			num := 1 + rand.Intn(10)
			events := make([]mq.Event, num)
			for i := 0; i < num; i++ {
				events[i] = createEvent()
			}

			docIdList, err := producer.Trigger(ctx, topic, events...)
			if err != nil {
				utils.Red("producer msg failed with error: %v", err)
			} else {
				utils.Green("producer msg list success and got id list: %v", docIdList)
			}

			time.Sleep(time.Second * 10)
		}()
	}
}
