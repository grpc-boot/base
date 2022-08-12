package main

import (
	"math"
	"math/rand"
	"runtime"
	"time"

	"github.com/grpc-boot/base/core/gopool"

	"github.com/grpc-boot/base"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	pool, err := base.NewGoPool(100,
		gopool.WithQueueLength(50),
		gopool.WithSpawnSize(runtime.NumCPU()),
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

	var (
		max = math.MaxInt16
	)

	for index := 0; index < 8; index++ {
		go func() {
			for i := 0; i < max; i++ {
				submit(pool)
			}
		}()
	}

	for {
		time.Sleep(time.Millisecond)
		if pool.PendingTaskTotal() == 0 {
			break
		}
	}

	base.Green("done workerNum:%d queueLen: %d pendingTaskTotal:%d successTotal:%d failedTotal:%d handleTotal:%d", pool.ActiveWorkerNum(), pool.QueueLength(), pool.PendingTaskTotal(), pool.SuccessTotal(), pool.FailedTotal(), pool.HandleTotal())
}

func submit(pool *gopool.Pool) {
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
