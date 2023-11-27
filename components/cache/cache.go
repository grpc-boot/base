package cache

import (
	"fmt"
	"hash/adler32"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/grpc-boot/base/v2/utils"

	"go.uber.org/atomic"
)

const (
	bucketLen   = 1<<9 - 1
	lockTimeout = 10 * time.Second
)

type Cache struct {
	data          [bucketLen]bucket
	length        atomic.Int64
	localDir      string
	latestSync    atomic.Time
	flushInterval time.Duration
}

func New(localDir string, flushInterval time.Duration) *Cache {
	if localDir != "" && !strings.HasSuffix(localDir, "/") {
		localDir += "/"
	}

	cache := &Cache{
		data:          [bucketLen]bucket{},
		localDir:      localDir,
		flushInterval: flushInterval,
	}

	cache.loadFromLocal()

	go cache.autoFlush()

	return cache
}

func (c *Cache) index(key string) int {
	return int(adler32.Checksum([]byte(key))) & bucketLen
}

func (c *Cache) GetValue(key string, timeoutSecond int64) (value []byte, exists bool) {
	var effective bool
	value, _, effective, exists = c.data[c.index(key)].getValue(key, timeoutSecond)
	if effective {
		return value, exists
	}
	return nil, exists
}

func (c *Cache) Get(key string, timeoutSecond int64, handler func() ([]byte, error)) []byte {
	index := c.index(key)

	value, lock, effective, exists := c.data[index].getValue(key, timeoutSecond)
	if !exists {
		// 第一次访问，初始化缓存
		data, err := handler()
		if err != nil {
			return nil
		}

		// 更新缓存
		c.SetValue(key, data)
		return data
	}

	// 缓存是否有效
	if effective {
		return value
	}

	// 缓存无效

	// 加锁
	token := utils.Acquire(lock, lockTimeout)
	// 加锁失败
	if token == 0 {
		return value
	}

	// 加锁成功，执行耗时操作
	data, err := handler()
	// 未获取到数据
	if err != nil {
		return nil
	}

	// 获取数据成功

	// 更新缓存
	c.SetValue(key, data)

	// 释放锁
	utils.Release(lock, token)

	return data
}

func (c *Cache) SetValue(key string, value []byte) {
	if isCreate := c.data[c.index(key)].setValue(key, value); isCreate {
		c.length.Inc()
	}
}

func (c *Cache) Delete(keys ...string) {
	for _, key := range keys {
		if num := c.data[c.index(key)].delete(key); num > 0 {
			c.length.Dec()
		}
	}
}

func (c *Cache) Length() int64 {
	return c.length.Load()
}

func (c *Cache) Info() Info {
	info := Info{
		LocalDir:       c.localDir,
		FlushInterval:  fmt.Sprintf("%s", c.flushInterval),
		LatestSyncTime: c.latestSync.Load().Format("2006-01-02 15:04:05"),
		Items:          make([]Item, 0, c.Length()),
	}

	for i := 0; i < bucketLen; i++ {
		info.Items = append(info.Items, c.data[i].items()...)
	}

	sort.SliceStable(info.Items, func(i, j int) bool {
		totalI := info.Items[i].Hit + info.Items[i].Miss
		totalJ := info.Items[j].Hit + info.Items[j].Miss
		if totalI > totalJ {
			return true
		}

		if info.Items[i].InvokeCount > info.Items[j].InvokeCount {
			return true
		}

		return info.Items[i].Miss > info.Items[j].Miss
	})

	info.Length = int64(len(info.Items))

	return info
}

func (c *Cache) SyncLocal() {
	if c.localDir == "" {
		return
	}

	for i := 0; i < bucketLen; i++ {
		_ = c.data[i].flushFile(c.localFileName(i))
	}
}

func (c *Cache) enableLocal() bool {
	return c.localDir != ""
}

func (c *Cache) loadFromLocal() {
	if !c.enableLocal() {
		return
	}

	for i := 0; i < bucketLen; i++ {
		length, _ := c.data[i].loadFile(c.localFileName(i))
		if length > 0 {
			c.length.Add(length)
		}
	}

	c.latestSync.Store(time.Now())
}

func (c *Cache) localFileName(index int) string {
	return c.localDir + "c-" + strconv.Itoa(index) + ".bin"
}

func (c *Cache) autoFlush() {
	if c.flushInterval < 1 || !c.enableLocal() {
		return
	}

	tick := time.NewTicker(c.flushInterval)
	for range tick.C {
		c.SyncLocal()
	}
}
