package cache

import (
	"hash/adler32"
	"strconv"
	"strings"
	"time"

	"go.uber.org/atomic"
)

const (
	bucketLen = 1<<9 - 1
)

type Cache struct {
	data          [bucketLen]bucket
	length        atomic.Int64
	localDir      string
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

func (c *Cache) Get(key string, timeoutSecond int64) (value []byte, exists bool) {
	var effective bool
	value, effective, exists = c.data[c.index(key)].getValue(key, timeoutSecond)
	if effective {
		return value, exists
	}
	return nil, exists
}

func (c *Cache) Common(key string, timeoutSecond int64, handler func() ([]byte, error)) []byte {
	index := c.index(key)

	item, exists := c.data[index].get(key)
	if !exists {
		// 第一次访问，初始化缓存
		data, err := handler()
		if err != nil {
			return nil
		}

		c.Set(key, data)
		return data
	}

	// 缓存是否有效
	if item.effective(timeoutSecond) {
		return item.Value
	}

	// 缓存无效，加锁
	token := item.lock()
	if token == 0 {
		return item.Value
	}

	// 加锁成功，执行耗时操作
	data, err := handler()
	if err != nil {
		return nil
	}

	item.addInvoke()
	item.CreatedAt = time.Now().Unix()
	item.Value = data

	//释放锁x
	item.unLock(token)
	return data
}

func (c *Cache) Set(key string, value []byte) {
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
