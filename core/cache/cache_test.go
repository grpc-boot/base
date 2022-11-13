package cache

import (
	jsoniter "github.com/json-iterator/go"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

var (
	localDir = "./data/"
)

func getKey(index int) string {
	return "c-key:" + strconv.Itoa(index)
}

func TestCache_Get(t *testing.T) {
	cache := New(localDir, time.Second)
	for i := 0; i < 1024; i++ {
		key := getKey(i)

		cache.SetValue(key, []byte(key))

		_, exists := cache.GetValue(getKey(i), 10)
		if !exists {
			t.Fatalf("want true, got %t", exists)
		}
	}

	info, _ := jsoniter.Marshal(cache.Info())
	t.Logf("cache info: %s", info)
}

func TestCache_SyncLocal(t *testing.T) {
	start := time.Now()
	runtime.GC()
	gcEnd := time.Now()

	t.Logf("first gc cost:%s", gcEnd.Sub(start))

	cache := New(localDir, time.Second)

	loadEnd := time.Now()
	t.Logf("load local cache cost:%s", loadEnd.Sub(gcEnd))

	for i := 0; i < 10000; i++ {
		key := getKey(i)

		cache.SetValue(key, []byte(strings.Repeat(key, 10)))
	}

	setEnd := time.Now()
	t.Logf("setValue 1000000 data cost:%s", setEnd.Sub(loadEnd))

	cache.SyncLocal()
	syncEnd := time.Now()

	t.Logf("sync local cost:%s", syncEnd.Sub(setEnd))
	runtime.GC()

	t.Logf("gc cost:%s", time.Now().Sub(syncEnd))
}

func TestCache_Common(t *testing.T) {
	cache := New(localDir, time.Second)
	maxCount := 16
	tick := time.NewTicker(time.Millisecond * 500)
	for range tick.C {
		data := cache.Get(getKey(1), 5, func() ([]byte, error) {
			time.Sleep(time.Second)
			return []byte(time.Now().String()), nil
		})
		t.Logf("data: %s", data)
		maxCount--
		if maxCount < 0 {
			break
		}
	}

	value, _ := cache.GetValue(getKey(1), 10)
	t.Logf("%s", value)
}

func TestCache_Commons(t *testing.T) {
	cache := New(localDir, time.Second)
	key := getKey(1)
	value, exists := cache.GetValue(key, 60)
	t.Logf("%s %t", value, exists)

	maxParallel := 8
	maxCount := 1000000

	wa := &sync.WaitGroup{}

	wa.Add(maxCount * maxParallel)

	for i := 0; i < maxParallel; i++ {
		go func() {
			for j := 0; j < maxCount; j++ {
				cache.Get(key, 1, func() ([]byte, error) {
					t.Logf("执行业务代码")
					time.Sleep(time.Second * 10)
					return []byte("asdfasdasdf"), nil
				})
				wa.Done()
			}
		}()
	}

	wa.Wait()

	cache.SyncLocal()
}
