package cache

import (
	"fmt"
	"sort"

	"github.com/grpc-boot/base/v2/kind/msg"
	"github.com/grpc-boot/base/v2/utils"
)

func Set[V msg.Value](c *Cache, key string, value V) (err error) {
	var isCreate bool
	isCreate, err = c.data[c.index(key)].set(key, value)
	if isCreate {
		c.length.Add(1)
	}
	return
}

func Get(c *Cache, key string, timeoutSecond int64) (value interface{}, effective, exists bool) {
	value, _, effective, exists = c.data[c.index(key)].get(key, timeoutSecond)
	return
}

func Delete(c *Cache, keys ...string) {
	for _, key := range keys {
		if num := c.data[c.index(key)].delete(key); num > 0 {
			c.length.Add(-1)
		}
	}
}

func Length(c *Cache) int64 {
	return c.length.Load()
}

func CommonGet[V msg.Value](c *Cache, key string, timeoutSecond int64, handler msg.Handler[V]) (interface{}, error) {
	var (
		hValue                       V
		err                          error
		index                        = c.index(key)
		val, lock, effective, exists = c.data[index].get(key, timeoutSecond)
	)

	if !exists {
		// 第一次访问，初始化缓存
		hValue, err = handler()
		if err != nil {
			return hValue, err
		}

		// 更新缓存
		err = Set(c, key, hValue)
		return hValue, err
	}

	// 缓存是否有效
	if effective {
		return val, nil
	}

	// 缓存无效

	// 加锁
	token := utils.Acquire(lock, lockTimeout)
	// 加锁失败
	if token == 0 {
		return val, nil
	}

	// 加锁成功，执行耗时操作
	hValue, err = handler()
	// 未获取到数据
	if err != nil {
		return hValue, err
	}

	// 获取数据成功

	// 更新缓存
	err = Set(c, key, hValue)
	return hValue, err
}

func InfoCache(c *Cache) Info {
	info := Info{
		LocalDir:       c.localDir,
		FlushInterval:  fmt.Sprintf("%s", c.flushInterval),
		LatestSyncTime: c.latestSync.Load().Format("2006-01-02 15:04:05"),
		Keys:           make(map[int][]string, bucketLen),
		Items:          make([]Item, 0, Length(c)+32),
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
