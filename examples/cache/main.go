package main

import (
	"time"

	"github.com/grpc-boot/base/v2/cache"
	"github.com/grpc-boot/base/v2/kind/msg"
	"github.com/grpc-boot/base/v2/utils"
)

func main() {
	localPath := "./cache"
	c := cache.New(localPath, time.Second*5)

	defer func() {
		// 手动同步数据到本地，运行过程中会自动同步
		c.Flush()
	}()

	value := cache.CommonGet[msg.Map](c, "index:conf", 10, func() (msg.Map, error) {
		// 模拟耗时
		time.Sleep(time.Second)

		return msg.Map{
			"rate":       3.14,
			"text":       "cache test",
			"updated_at": time.Now().Unix(),
		}, nil
	})

	conf := msg.MsgMap(value)
	utils.Green("rate: %.2f text: %s updated at: %d", conf.Float("rate"), conf.String("text"), conf.Int("updated_at"))
}
