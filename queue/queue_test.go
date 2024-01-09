package queue

import (
	"context"
	"golang.org/x/exp/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/grpc-boot/base/v2/gored"
)

func init() {
	o := gored.DefaultOptions()
	gored.SetRedis("redis", o)
}

func TestDelayRedis_Fetch(t *testing.T) {
	var (
		red, _   = gored.GetRedis("redis")
		opt      = DefaultOptions()
		dq       = NewDelay("delay_queue_test", red, opt)
		maxCount = 10000
		wa       sync.WaitGroup
	)

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

	for i := 0; i < maxCount; i++ {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			item := Item{
				Id:   time.Now().Format("20060102150405") + strconv.Itoa(rand.Int()),
				Name: "remove order",
				At:   time.Now().Add(time.Second * 10).Unix(),
			}

			err = dq.Set(ctx, item)
			if err != nil {
				t.Fatalf("want nil, got %v", err)
			}
		}()
	}

	wa.Wait()
}
