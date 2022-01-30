package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/grpc-boot/base"
	"github.com/grpc-boot/base/core/shardmap"
)

var (
	done chan struct{}
)

func main() {
	sm, ch := base.NewSharMapWithChannel(1024)

	// 模拟消费
	go func() {
		for {
			event := <-ch
			switch event.Type {
			case shardmap.Create:
				base.Green("create key:%+v value:%+v", event.Key, event.Value)
			case shardmap.Update:
				base.Green("update key:%+v value:%+v oldValue:%+v", event.Key, event.Value, event.OldValue)
			case shardmap.Delete:
				base.Green("delete key:%+v", event.Key)
			}
		}
	}()

	// 模拟写入
	go func() {
		for {
			rand.Seed(time.Now().UnixNano())
			opt := rand.Intn(3)
			keyFormat := "key%d"
			switch opt {
			case 0:
				sm.Set(fmt.Sprintf(keyFormat, 0), time.Now().UnixMicro())
			case 1:
				sm.Set(fmt.Sprintf(keyFormat, 1), time.Now().UnixMicro())
			case 2:
				sm.Delete(fmt.Sprintf(keyFormat, rand.Intn(3)))
			}

			time.Sleep(time.Millisecond * 300)
		}
	}()

	<-done
}
