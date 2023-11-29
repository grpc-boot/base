package utils

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"go.uber.org/atomic"
)

func TestAcquire_WithoutRelease(t *testing.T) {
	var (
		workerNum = 128
		waitTime  = time.Second * 30
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

				token := Acquire(&locker, time.Second)
				if token > 0 {
					t.Logf("%d: got token:%d", time.Now().Unix(), token)
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
