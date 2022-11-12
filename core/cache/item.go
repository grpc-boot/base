//go:generate msgp
package cache

import (
	"sync/atomic"
	"time"
)

const (
	//单位纳秒
	lockTimeout = 10 * 1000 * 1000 * 1000
)

type Bucket struct {
	Items map[string]*Item `json:"items" msg:"items"`
}

type Item struct {
	_lock       int64
	Value       []byte `json:"value" msg:"value"`
	Hit         int64  `json:"hit" msg:"hit"`
	Miss        int64  `json:"miss" msg:"miss"`
	InvokeCount int64  `json:"invokeCount" msg:"invokeCount"`
	CreatedAt   int64  `json:"createdAt" msg:"createdAt"`
}

// effective 数据是否有效
func (i *Item) effective(timeoutSecond int64) (ok bool) {
	ok = i.CreatedAt+timeoutSecond > time.Now().Unix()
	if ok {
		atomic.AddInt64(&i.Hit, 1)
	} else {
		atomic.AddInt64(&i.Miss, 1)
	}

	return ok
}

func (i *Item) addInvoke() {
	atomic.AddInt64(&i.InvokeCount, 1)
}

// lock 加锁
func (i *Item) lock() int64 {
	now := time.Now().UnixNano()

	// 尝试获取锁
	if atomic.CompareAndSwapInt64(&i._lock, 0, now) {
		return now
	}

	current := atomic.LoadInt64(&i._lock)
	// 加锁时间已经超过lockTimeout
	if now-current > lockTimeout {
		// 重新加锁
		if atomic.CompareAndSwapInt64(&i._lock, current, now) {
			return now
		}
	}
	return 0
}

// unLock 释放自己加的锁
func (i *Item) unLock(token int64) {
	atomic.CompareAndSwapInt64(&i._lock, token, 0)
}
