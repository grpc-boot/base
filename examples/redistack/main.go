package main

import (
	"time"

	"github.com/grpc-boot/base/v2/redistack"
	"github.com/grpc-boot/base/v2/utils"

	"github.com/redis/rueidis"
)

var (
	client    rueidis.Client
	streamKey = `key:stream`
	group     = `group2`
	consumer  = `consumer2`
)

func init() {
	var err error
	client, err = rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{"10.16.49.131:6379"},
	})

	if err != nil {
		panic(err)
	}
}

func main() {
	var (
		// 故障恢复
		id     = `0`
		msgMap map[string][]rueidis.XRangeEntry
	)

	for {
		info, err := redistack.StreamInfo(time.Second, client, streamKey)
		if err != nil {
			utils.RedFatal("fetch stream info failed with error: %v", err)
		}

		utils.Green("fetch stream info: %+v", info)

		msgMap, err = redistack.StreamReadGroup(10*1000, client, streamKey, group, consumer, 10, id, false)
		if !redistack.IsNil(err) {
			utils.RedFatal("consumer msg failed with error: %v", err)
		}

		if _, exists := msgMap[streamKey]; exists && len(msgMap[streamKey]) == 0 && id == "0" {
			//获取所有未传递的消息
			id = `>`
		}

		for channel, list := range msgMap {
			if len(list) == 0 {
				continue
			}

			idList := make([]string, len(list))
			for index, msg := range list {
				utils.Green("[%s] got msg:[%s] %+v", channel, msg.ID, msg.FieldValues)
				idList[index] = msg.ID
			}
			okCount, err := redistack.StreamAck(time.Second*3, client, streamKey, group, idList...)
			if err != nil {
				utils.RedFatal("ack failed with error: %v", err)
			}
			utils.Green("ack success count: %d", okCount)
		}
	}
}
