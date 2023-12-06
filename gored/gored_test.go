package gored

import (
	"context"
	"testing"
	"time"

	"github.com/grpc-boot/base/v2/kind/msg"

	"github.com/redis/go-redis/v9"
)

var (
	opt = func() redis.Options {
		o := DefaultOptions()
		o.Addr = "127.0.0.1:6379"
		return o
	}()
)

func TestGetItemWithCache(t *testing.T) {
	red := redis.NewClient(&opt)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	item, err := GetItemWithCache(ctx, red, "cache", time.Now().Unix(), 6, func() (value msg.Map, err error) {
		value = msg.Map{
			"id":   10086,
			"name": "移动",
		}
		return
	})

	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	t.Logf("value: %v", item.Map())
}

func TestAcquire(t *testing.T) {
	red := redis.NewClient(&opt)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	cmd := Acquire(ctx, red, "acquire", 10)
	err := DealCmdErr(cmd)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	token := cmd.Val()
	if token > 0 {
		t.Logf("acquire token: %d", token)
		ctx1, cancel1 := context.WithTimeout(context.Background(), time.Second)
		defer cancel1()

		rCmd := Release(ctx1, red, "acquire", token)
		err = DealCmdErr(rCmd)
		if err != nil {
			t.Fatalf("want nil, got %v", err)
		}
	} else {
		t.Logf("do not acquire token")
	}
}

func TestAcquireWithTimeout(t *testing.T) {
	red := redis.NewClient(&opt)

	cmd := AcquireWithTimeout(time.Second, red, "acquire", 10)
	err := DealCmdErr(cmd)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	token := cmd.Val()
	if token > 0 {
		t.Logf("acquire token: %d", token)

		rCmd := ReleaseWithTimeout(time.Second, red, "acquire", token)
		err = DealCmdErr(rCmd)
		if err != nil {
			t.Fatalf("want nil, got %v", err)
		}
	} else {
		t.Logf("do not acquire token")
	}
}
