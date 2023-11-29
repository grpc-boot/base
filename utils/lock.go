package utils

import (
	"sync/atomic"
	"time"
)

func Acquire(lock *int64, timeout time.Duration) (token int64) {
	now := time.Now().UnixNano()
	if atomic.CompareAndSwapInt64(lock, 0, now) {
		return now
	}

	lastToken := atomic.LoadInt64(lock)
	if time.Duration(now-lastToken) > timeout {
		if atomic.CompareAndSwapInt64(lock, lastToken, now) {
			return now
		}
	}

	return 0
}

func Release(lock *int64, token int64) (swapped bool) {
	return atomic.CompareAndSwapInt64(lock, token, 0)
}
