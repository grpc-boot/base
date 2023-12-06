package cache

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/grpc-boot/base/v2/kind"
	"github.com/grpc-boot/base/v2/utils"

	"go.uber.org/atomic"
)

const (
	bucketLen   = 1 << 8
	lockTimeout = 10 * time.Second
)

var (
	FlushWithoutChangeIntervalSeconds int64 = 60
)

type Cache struct {
	data          [bucketLen]Bucket
	length        atomic.Int64
	localDir      string
	latestSync    atomic.Time
	flushInterval time.Duration
}

func New(localDir string, flushInterval time.Duration) *Cache {
	if localDir != "" && localDir[len(localDir)-1] != '/' {
		localDir += "/"
	}

	if localDir != "" {
		_ = utils.MkDir(localDir, 0766)
	}

	cache := &Cache{
		data:          [bucketLen]Bucket{},
		localDir:      localDir,
		flushInterval: flushInterval,
	}

	cache.loadFromLocal()

	go cache.autoFlush()

	return cache
}

func (c *Cache) index(key string) int {
	return int(kind.Uint32Hash(utils.String2Bytes(key))) & (bucketLen - 1)
}

func (c *Cache) SyncLocal() {
	if c.localDir == "" {
		return
	}

	for i := 0; i < bucketLen; i++ {
		utils.Recover(fmt.Sprintf("sync cache[%d] to local file", i), func(args ...any) {
			_ = c.data[i].flushFile(c.localFileName(i))
		})
	}
}

func (c *Cache) enableLocal() bool {
	return c.localDir != ""
}

func (c *Cache) loadFromLocal() {
	if !c.enableLocal() {
		return
	}

	var wg sync.WaitGroup
	wg.Add(bucketLen)

	for i := 0; i < bucketLen; i++ {
		go utils.Recover("", func(args ...any) {
			defer func() {
				args[0].(*sync.WaitGroup).Done()
			}()

			index := args[1].(int)
			length, _ := c.data[index].loadFile(c.localFileName(index))
			if length > 0 {
				c.length.Add(length)
			}
		}, &wg, i)
	}

	wg.Wait()

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
