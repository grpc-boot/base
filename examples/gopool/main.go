package main

import (
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/grpc-boot/base/core/gopool"

	"github.com/grpc-boot/base"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	pool, err := base.NewGoPool(100,
		//gopool.WithQueueLength(50),
		//gopool.WithSpawnSize(10),
		gopool.WithPanicHandler(func(err interface{}) {
			if err != nil {
				base.Red("panic error:%+v\n", err)
			}
		}),
	)

	if err != nil {
		base.RedFatal("new pool error:%s\n", err.Error())
	}

	go func() {
		base.Green("workerNum:%d queueLen: %d", pool.ActiveWorkerNum(), pool.QueueLength())

		tick := time.NewTicker(time.Second)
		for range tick.C {
			base.Green("workerNum:%d queueLen: %d", pool.ActiveWorkerNum(), pool.QueueLength())
		}
	}()

	var (
		wg  sync.WaitGroup
		max = math.MaxUint16
	)

	wg.Add(max)

	for index := 0; index < max; index++ {
		err = pool.Submit(func() {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			wg.Done()
		})

		if err != nil {
			wg.Done()
			base.Red("submit error:%s\n", err.Error())
		}
	}

	wg.Wait()
	base.Green("done workerNum:%d queueLen: %d", pool.ActiveWorkerNum(), pool.QueueLength())
}
