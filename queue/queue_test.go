package queue

import (
	"context"
	"fmt"
	"golang.org/x/exp/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/grpc-boot/base/v2/gopool"
	"github.com/grpc-boot/base/v2/gored"
	"github.com/grpc-boot/base/v2/utils"
)

var (
	goPool *gopool.Pool
)

func init() {

	o := gored.DefaultOptions()
	gored.SetRedis("redis", o)

	/*var err error
	goPool, err = gopool.NewPool(8,
		gopool.WithQueueLength(64),
		gopool.WithSpawnSize(1),
		gopool.WithMaxIdleTimeoutSeconds(30),
		gopool.WithPanicHandler(func(err interface{}) {
			if err != nil {
				utils.Red("panic error:%+v", err)
			}
		}),
	)

	if err != nil {
		panic(err)
	}*/
}

func TestDelayRedis_Fetch(t *testing.T) {
	var (
		red, _   = gored.GetRedis("redis")
		opt      = DefaultOptions()
		maxCount = 1000
		wa       sync.WaitGroup
	)

	dq := NewDelay("delay_queue_test", red, opt)

	wa.Add(maxCount)

	defer func() {
		t.Logf("call stop at: %s", time.Now().Format(time.DateTime))
		err := dq.Stop(time.Second * 10)
		if err != nil {
			t.Fatalf("stop with error:%v", err)
		}

		t.Logf("stop success at: %s", time.Now().Format(time.DateTime))
	}()

	dq.RegisterHandler(func(items []Item) {
		t.Logf("got count: %d\n", len(items))

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		_ = dq.Done(ctx, items...)

		for i := 0; i < len(items); i++ {
			if items[i].RetryCount == 0 && i%2 == 0 {
				continue
			}

			wa.Done()
		}
	})

	err := dq.Start()
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	go func() {
		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			utils.Green("info tick")
			func() {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				list, er := dq.DeadList(ctx)
				if er != nil {
					fmt.Printf("fetch dead list failed with error:%v\n", er)
					return
				}
				utils.Green("fetch dead list: %+v", list)
				/*utils.Green("workerNum:%d queueLen: %d pendingTaskTotal:%d successTotal:%d failedTotal:%d handleTotal:%d", goPool.ActiveWorkerNum(), goPool.QueueLength(), goPool.PendingTaskTotal(), goPool.SuccessTotal(), goPool.FailedTotal(), goPool.HandleTotal())*/
			}()
		}
	}()

	for i := 0; i < maxCount; i++ {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			item := Item{
				Id:   time.Now().Format("20060102150405") + strconv.Itoa(rand.Int()),
				Name: fmt.Sprintf("remove order: %d", time.Now().UnixNano()),
				At:   time.Now().Add(time.Second * time.Duration(1+rand.Int63n(1000))).Unix(),
			}

			err = dq.Set(ctx, item)
			if err != nil {
				t.Fatalf("want nil, got %v", err)
			}
		}()
	}

	utils.Green("set done: %d", time.Now().Unix())

	wa.Wait()

	utils.Green("end at: %d", time.Now().Unix())
}
