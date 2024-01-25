# base

### cache 本地内存缓存组件，会同步到本地目录，重启后从本地加载到内存

```go
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
```

### `gored`操作redis

> 缓存处理

```go
package gored

import (
	"context"
	"testing"
	"time"

	"github.com/grpc-boot/base/v2/kind/msg"
)

func init() {
	o := DefaultOptions()
	SetRedis("redis", o)
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
```

> 锁

```go
func TestAcquire(t *testing.T) {
	red, _ := GetRedis("redis")

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
```

> 令牌桶限速

```go
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
```

### 运行时开启关闭pprof，当系统出现问题时可以实时开启pprof定位系统问题

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-boot/base/v2/internal"
	"github.com/grpc-boot/base/v2/utils"
)

type router struct {
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	_, _ = w.Write(internal.String2Bytes(`ok`))
}

func main() {
	server := http.Server{
		Addr:    ":8080",
		Handler: &router{},
	}

	go func() {
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			panic(err)
		}
	}()

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		err := server.Shutdown(ctx)
		if err != nil {
			fmt.Printf("shutdown server error:%v", err)
		}
	}()

	var sig = make(chan os.Signal, 1)
	signal.Notify(sig)

	for {
		val := <-sig
		switch val {
		case syscall.SIGUSR1:
			if utils.PprofIsRun() {
				continue
			}

			go func() {
				err := utils.StartPprof(":8081", nil)
				if err != nil {
					fmt.Printf("start pprof error:%v", err)
				}
			}()
		case syscall.SIGUSR2:
			if !utils.PprofIsRun() {
				continue
			}

			err := utils.StopPprofWithTimeout(10)
			if err != nil {
				fmt.Printf("stop pprof error:%v", err)
			}
		default:
			signal.Stop(sig)
			return
		}
	}
}
```

```shell
## 开启pprof
kill -USR1 ${pid}

## 关闭pprof
kill -USR2 ${pid}
```

### utils.Timeout帮助方法，简化超时实现逻辑

```go
func TestTimeout(t *testing.T) {
    err := Timeout(time.Second, func(args ...any) {
        time.Sleep(time.Millisecond * 500)
    })
    
    if err != nil {
        t.Fatalf("want nil, got %v", err)
    }
    
    err = Timeout(time.Millisecond*100, func(args ...any) {
        time.Sleep(time.Millisecond * 200)
    })
    
    if err != context.DeadlineExceeded {
        t.Fatalf("want err, got %v", err)
    }
}
```

### `utils.Recover`帮助方法，减少未知panic导致进程宕掉

```go
func TestRecover(t *testing.T) {
	go Recover("recover test", func(args ...any) {
		panic("panic with test")
	})
}
```

### utils.Join 可以对int、int32等数值类型进行Join

```go
func TestJoin(t *testing.T) {
	ss := []string{"s1", "s2"}

	res1 := strings.Join(ss, ",")
	t.Logf("res1: %s", res1)

	is := []int{1, 2, 45}
	resInt := Join(",", is...)
	t.Logf("resInt: %s", resInt)

	i32s := []int32{1, 2, 45}
	resInt32 := Join(",", i32s...)
	t.Logf("resInt32: %s", resInt32)
}
```

### utils.Acquire基于原子操作的超时锁

```go
func TestAcquire(t *testing.T) {
	var (
		workerNum = 32
		waitTime  = time.Second * 60
		done      atomic.Bool
		locker    int64
		wa        sync.WaitGroup
	)

	wa.Add(workerNum)

	for i := 0; i < workerNum; i++ {
		go func(w *sync.WaitGroup) {
			for {
				if done.Load() {
					break
				}

				token := Acquire(&locker, time.Second*3)
				if token > 0 {
					// 模拟业务操作成功，释放锁，操作失败，保留锁，防止频繁获得锁操作业务
					ok := func() bool {
						// load data from db to cache

						return rand.Int63()%2 == 0
					}()

					if ok {
						Release(&locker, token)
					}
				}
			}
			w.Done()
		}(&wa)
	}

	time.AfterFunc(waitTime, func() {
		done.Store(true)
	})

	wa.Wait()
}
```


