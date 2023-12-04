package cache

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/grpc-boot/base/v2/kind"
	"github.com/grpc-boot/base/v2/utils"

	"go.uber.org/atomic"
)

const (
	bucketLen   = 1<<7 - 1
	lockTimeout = 10 * time.Second
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
	return int(kind.Uint32Hash(utils.String2Bytes(key))) & bucketLen
}

func (c *Cache) Set(key string, value interface{}) (err error) {
	var isCreate bool
	isCreate, err = c.data[c.index(key)].set(key, value)
	if isCreate {
		c.length.Add(1)
	}
	return
}

func (c *Cache) Get(key string, timeoutSecond int64) (value interface{}, effective, exists bool) {
	value, _, effective, exists = c.data[c.index(key)].get(key, timeoutSecond)
	return
}

func (c *Cache) GetBool(key string, timeoutSecond int64) (value bool, effective bool) {
	var (
		val    interface{}
		exists bool
	)

	val, effective, exists = c.Get(key, timeoutSecond)
	if !exists {
		return
	}

	value, _ = val.(bool)
	return
}

func (c *Cache) GetInt(key string, timeoutSecond int64) (value int64, effective bool) {
	var (
		val    interface{}
		exists bool
	)

	val, effective, exists = c.Get(key, timeoutSecond)
	if !exists {
		return
	}

	value = Int(val)
	return
}

func (c *Cache) GetUint(key string, timeoutSecond int64) (value uint64, effective bool) {
	var (
		val    interface{}
		exists bool
	)

	val, effective, exists = c.Get(key, timeoutSecond)
	if !exists {
		return
	}
	value = Uint(val)
	return
}

func (c *Cache) GetFloat(key string, timeoutSecond int64) (value float64, effective bool) {
	var (
		val    interface{}
		exists bool
	)

	val, effective, exists = c.Get(key, timeoutSecond)
	if !exists {
		return
	}
	value = Float(val)
	return
}

func (c *Cache) GetString(key string, timeoutSecond int64) (value string, effective bool) {
	var (
		val    interface{}
		exists bool
	)

	val, effective, exists = c.Get(key, timeoutSecond)
	if !exists {
		return
	}

	value, _ = val.(string)
	return
}

func (c *Cache) GetBytes(key string, timeoutSecond int64) (value []byte, effective bool) {
	var (
		val    interface{}
		exists bool
	)

	val, effective, exists = c.Get(key, timeoutSecond)
	if !exists {
		return
	}

	value, _ = val.([]byte)
	return
}

func (c *Cache) GetSlice(key string, timeoutSecond int64) (value []interface{}, effective bool) {
	var (
		val    interface{}
		exists bool
	)

	val, effective, exists = c.Get(key, timeoutSecond)
	if !exists {
		return
	}

	value, _ = val.([]interface{})
	return
}

func (c *Cache) GetMap(key string, timeoutSecond int64) (value Map, effective bool) {
	var (
		val    interface{}
		exists bool
	)

	val, effective, exists = c.Get(key, timeoutSecond)
	if !exists {
		return
	}

	v, _ := val.(map[string]interface{})
	value = Map(v)
	return
}

func (c *Cache) GetTime(key string, timeoutSecond int64) (value time.Time, effective bool) {
	var (
		val    interface{}
		exists bool
	)

	val, effective, exists = c.Get(key, timeoutSecond)
	if !exists {
		return
	}

	value, _ = val.(time.Time)
	return
}

func (c *Cache) Delete(keys ...string) {
	for _, key := range keys {
		if num := c.data[c.index(key)].delete(key); num > 0 {
			c.length.Add(-1)
		}
	}
}

func (c *Cache) CommonGet(key string, timeoutSecond int64, handler func() (interface{}, error)) (value interface{}, err error) {
	var (
		index                        = c.index(key)
		val, lock, effective, exists = c.data[index].get(key, timeoutSecond)
	)

	if !exists {
		// 第一次访问，初始化缓存
		value, err = handler()
		if err != nil {
			return
		}

		// 更新缓存
		err = c.Set(key, value)
		return
	}

	// 缓存是否有效
	if effective {
		value = val
		return
	}

	// 缓存无效

	// 加锁
	token := utils.Acquire(lock, lockTimeout)
	// 加锁失败
	if token == 0 {
		value = val
		return
	}

	// 加锁成功，执行耗时操作
	data, err := handler()
	// 未获取到数据
	if err != nil {
		return
	}

	// 获取数据成功

	// 更新缓存
	err = c.Set(key, data)
	return
}

func (c *Cache) CommonBool(key string, timeoutSecond int64, handler func() (interface{}, error)) (value bool, err error) {
	val, err := c.CommonGet(key, timeoutSecond, handler)
	if err != nil {
		return
	}

	value, _ = val.(bool)
	return
}

func (c *Cache) CommonInt(key string, timeoutSecond int64, handler func() (interface{}, error)) (value int64, err error) {
	val, err := c.CommonGet(key, timeoutSecond, handler)
	if err != nil {
		return
	}

	value = Int(val)
	return
}

func (c *Cache) CommonFloat(key string, timeoutSecond int64, handler func() (interface{}, error)) (value float64, err error) {
	val, err := c.CommonGet(key, timeoutSecond, handler)
	if err != nil {
		return
	}

	value = Float(val)
	return
}

func (c *Cache) CommonUint(key string, timeoutSecond int64, handler func() (interface{}, error)) (value uint64, err error) {
	val, err := c.CommonGet(key, timeoutSecond, handler)
	if err != nil {
		return
	}

	value = Uint(val)
	return
}

func (c *Cache) CommonString(key string, timeoutSecond int64, handler func() (interface{}, error)) (value string, err error) {
	val, err := c.CommonGet(key, timeoutSecond, handler)
	if err != nil {
		return
	}

	value, _ = val.(string)
	return
}

func (c *Cache) CommonBytes(key string, timeoutSecond int64, handler func() (interface{}, error)) (value []byte, err error) {
	val, err := c.CommonGet(key, timeoutSecond, handler)
	if err != nil {
		return
	}

	value, _ = val.([]byte)
	return
}

func (c *Cache) CommonSlice(key string, timeoutSecond int64, handler func() (interface{}, error)) (value []interface{}, err error) {
	val, err := c.CommonGet(key, timeoutSecond, handler)
	if err != nil {
		return
	}

	value, _ = val.([]interface{})
	return
}

func (c *Cache) CommonMap(key string, timeoutSecond int64, handler func() (interface{}, error)) (value Map, err error) {
	val, err := c.CommonGet(key, timeoutSecond, handler)
	if err != nil {
		return
	}

	v, _ := val.(map[string]interface{})
	value = Map(v)
	return
}

func (c *Cache) Length() int64 {
	return c.length.Load()
}

func (c *Cache) Info() Info {
	info := Info{
		LocalDir:       c.localDir,
		FlushInterval:  fmt.Sprintf("%s", c.flushInterval),
		LatestSyncTime: c.latestSync.Load().Format("2006-01-02 15:04:05"),
		Keys:           make(map[int][]string, bucketLen),
		Items:          make([]Item, 0, c.Length()+32),
	}

	for i := 0; i < bucketLen; i++ {
		var (
			items = c.data[i].items()
			keys  = make([]string, len(items))
		)

		if len(items) > 0 {
			info.Items = append(info.Items, items...)

			for index, item := range items {
				keys[index] = item.Key
			}
			info.Keys[i] = keys
		}
	}

	sort.SliceStable(info.Items, func(i, j int) bool {
		totalI := info.Items[i].HitCnt + info.Items[i].MissCnt
		totalJ := info.Items[j].HitCnt + info.Items[j].MissCnt
		if totalI > totalJ {
			return true
		}

		if info.Items[i].UpdateCnt > info.Items[j].UpdateCnt {
			return true
		}

		return info.Items[i].MissCnt > info.Items[j].MissCnt
	})

	info.Length = int64(len(info.Items))

	return info
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
