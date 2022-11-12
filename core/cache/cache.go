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
	dir           string
	flushInterval time.Duration
}

func New(dir string) *Cache {
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}

	cache := &Cache{
		data:          [bucketLen]bucket{},
		dir:           dir,
		flushInterval: time.Second,
	}

	go cache.autoFlush()

	return cache
}

func (c *Cache) index(key string) int {
	return int(adler32.Checksum([]byte(key))) & bucketLen
}

func (c *Cache) Get(key string) (item *Item, exists bool) {
	return c.data[c.index(key)].get(key)
}

func (c *Cache) Set(key string, value []byte, timeoutSecond int64) {
	item := &Item{
		ExpireAt: time.Now().Unix() + timeoutSecond,
		Value:    value,
	}

	if isCreate := c.data[c.index(key)].set(key, item); isCreate {
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

func (c *Cache) autoFlush() {
	if c.flushInterval < 1 {
		return
	}

	tick := time.NewTicker(c.flushInterval)
	for range tick.C {
		c.flushFile()
	}
}

func (c *Cache) flushFile() {
	for i := 0; i < bucketLen; i++ {
		fileName := c.dir + "c-" + strconv.Itoa(i) + ".bin"
		_ = c.data[i].flushFile(fileName)
	}
}
