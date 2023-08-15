package main

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/grpc-boot/base"
	"github.com/grpc-boot/base/core/gopool"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	pool, err := base.NewGoPool(100,
		gopool.WithQueueLength(50),
		gopool.WithSpawnSize(runtime.NumCPU()),
		gopool.WithMaxIdleTimeoutSeconds(30),
		gopool.WithPanicHandler(func(err interface{}) {
			if err != nil {
				base.Red("panic error:%+v", err)
			}
		}),
	)

	if err != nil {
		base.RedFatal("new pool error:%s", err.Error())
	}

	go func() {
		base.Green("workerNum:%d queueLen: %d pendingTaskTotal:%d successTotal:%d failedTotal:%d handleTotal:%d", pool.ActiveWorkerNum(), pool.QueueLength(), pool.PendingTaskTotal(), pool.SuccessTotal(), pool.FailedTotal(), pool.HandleTotal())

		tick := time.NewTicker(time.Second)
		for range tick.C {
			base.Green("workerNum:%d queueLen: %d pendingTaskTotal:%d successTotal:%d failedTotal:%d handleTotal:%d", pool.ActiveWorkerNum(), pool.QueueLength(), pool.PendingTaskTotal(), pool.SuccessTotal(), pool.FailedTotal(), pool.HandleTotal())
		}
	}()

	for index := 0; index < 8; index++ {
		go submit(pool)
	}

	done := make(chan struct{}, 1)
	<-done

	base.Green("done workerNum:%d queueLen: %d pendingTaskTotal:%d successTotal:%d failedTotal:%d handleTotal:%d", pool.ActiveWorkerNum(), pool.QueueLength(), pool.PendingTaskTotal(), pool.SuccessTotal(), pool.FailedTotal(), pool.HandleTotal())
}

func submit(pool *gopool.Pool) {
	var (
		interval    = time.Millisecond * 10
		maxInterval = time.Second * 3
		ticker      = time.NewTicker(interval)
		resetTicker = time.NewTicker(maxInterval)
	)

	go func() {
		for range resetTicker.C {
			if interval < maxInterval {
				interval += time.Millisecond * 10
				ticker.Reset(interval)
			}
		}
	}()

	for range ticker.C {
		err := pool.Submit(func() {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(10)))
			val := rand.Intn(10000)
			if val < 1 {
				panic("hit")
			}
		})

		if err != nil {
			base.Red("submit error:%s", err.Error())
		}
	}
}
