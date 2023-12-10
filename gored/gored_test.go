package gored

import (
	"context"
	"testing"
	"time"

	"github.com/grpc-boot/base/v2/internal"
	"github.com/grpc-boot/base/v2/kind/msg"

	"github.com/redis/go-redis/v9"
)

func init() {
	o := DefaultOptions()
	SetRedis("redis", o)
}

func TestConf_Float(t *testing.T) {
	var (
		start    = time.Now()
		red, _   = GetRedis("redis")
		interval = time.Second * 3
		conf     = NewConf("base_conf", interval, red)
		dp       = map[string]interface{}{
			"int":    56,
			"int64":  int64(2313),
			"bool":   true,
			"string": `的as2`,
			"[]byte": []byte("sd都是"),
			"float":  float32(3434.32),
		}
	)

	t.Logf("load conf cost: %s", time.Since(start))

	_, err := SetConf(conf, "int", dp["int"].(int))
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	_, err = SetConf(conf, "int64", dp["int64"].(int64))
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	_, err = SetConf(conf, "bool", dp["bool"].(bool))
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	_, err = SetConf(conf, "string", dp["string"].(string))
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	_, err = SetConf(conf, "[]byte", dp["[]byte"].([]byte))
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	_, err = SetConf(conf, "float", dp["float"].(float32))
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	time.Sleep(interval + time.Second)

	time.Sleep(time.Second)

	intVal := conf.Int("int", 1)
	if intVal != int64(dp["int"].(int)) {
		t.Fatalf("want %d, got %v", dp["int"].(int), intVal)
	}

	int64Val := conf.Int("int64", 1)
	if int64Val != dp["int64"].(int64) {
		t.Fatalf("want %d, got %v", dp["int64"].(int64), int64Val)
	}

	boolVal := conf.Bool("bool", false)
	if boolVal != dp["bool"].(bool) {
		t.Fatalf("want %t, got %v", dp["bool"].(bool), boolVal)
	}

	stringVal := conf.String("string", "f")
	if stringVal != dp["string"].(string) {
		t.Fatalf("want %s, got %v", dp["string"].(string), stringVal)
	}

	floatVal := conf.Float("float", 3.1415)
	if floatVal != float64(dp["float"].(float32)) {
		t.Fatalf("want %v, got %v", dp["float"].(float32), floatVal)
	}

	bytesVal := conf.String("[]byte", "sadfas")
	if bytesVal != internal.Bytes2String(dp["[]byte"].([]byte)) {
		t.Fatalf("want %s, got %v", dp["[]byte"].([]byte), bytesVal)
	}
}

func TestGetItemWithCache(t *testing.T) {
	red, _ := GetRedis("redis")
	item, err := GetItemWithCacheTimeout(time.Second, red, "cache", time.Now().Unix(), 6, func() (value msg.Map, err error) {
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

// BenchmarkGetItemWithCache-8   	   21046	     52605 ns/op
func BenchmarkGetItemWithCache(b *testing.B) {
	red, _ := GetRedis("redis")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		item, err := GetItemWithCacheTimeout(time.Second, red, "cache", time.Now().Unix(), 6, func() (value msg.Map, err error) {
			value = msg.Map{
				"id":   10086,
				"name": "移动",
			}
			return
		})

		if err != nil {
			b.Fatalf("want nil, got %v", err)
		}

		id := item.Map().Int("id")
		if id != 10086 {
			b.Fatalf("want 10086, got %v", id)
		}
	}
}

// BenchmarkGetItemWithCacheParallel-8   	   18128	     77644 ns/op
func BenchmarkGetItemWithCacheParallel(b *testing.B) {
	red, _ := GetRedis("redis")

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			item, err := GetItemWithCacheTimeout(time.Second, red, "cache", time.Now().Unix(), 6, func() (value msg.Map, err error) {
				value = msg.Map{
					"id":   10086,
					"name": "移动",
				}
				return
			})

			if err != nil {
				b.Fatalf("want nil, got %v", err)
			}

			id := item.Map().Int("id")
			if id != 10086 {
				b.Fatalf("want 10086, got %v", id)
			}
		}
	})
}

func TestAcquire(t *testing.T) {
	var (
		red, _ = GetRedis("redis")
		cmd    *redis.IntCmd
	)

	TimeoutDo(time.Second, func(ctx context.Context) {
		cmd = Acquire(ctx, red, "acquire", 10)
	})

	err := DealCmdErr(cmd)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	token := cmd.Val()
	if token > 0 {
		t.Logf("acquire token: %d", token)

		var rCmd *redis.IntCmd
		TimeoutDo(time.Second, func(ctx context.Context) {
			rCmd = Release(ctx, red, "acquire", token)
		})

		err = DealCmdErr(rCmd)
		if err != nil {
			t.Fatalf("want nil, got %v", err)
		}
	} else {
		t.Logf("do not acquire token")
	}
}

func TestAcquireWithTimeout(t *testing.T) {
	red, _ := GetRedis("redis")

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

func TestGetToken(t *testing.T) {
	red, _ := GetRedis("redis")
	max := 100

	for i := 0; i < max; i++ {
		cmd := SecondLimitByToken(context.Background(), red, "token", 3, 1, 6)
		err := DealCmdErr(cmd)
		if err != nil {
			t.Fatalf("want nil, got %v", err)
		}

		if cmd.Val() {
			t.Logf("got token")
		}
	}
}

// BenchmarkSecondLimitByToken-8   	   19594	     56544 ns/op
func BenchmarkSecondLimitByToken(b *testing.B) {
	red, _ := GetRedis("redis")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cmd := SecondLimitByToken(context.Background(), red, "token", 3, 1, 6)
		err := DealCmdErr(cmd)
		if err != nil {
			b.Fatalf("want nil, got %v", err)
		}
	}
}
