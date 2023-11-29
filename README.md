# base

### utils.Recover帮助方法，减少未知panic导致进程宕掉

```go
func TestRecover(t *testing.T) {
	go Recover("recover test", func() {
		panic("panic with test")
	})
}
```

### utils.Join 可以对int、int32等数值类型进行Join

```go
func TestJoin(t *testing.T) {
	ss := []string{"s1", "s2"}

	res1 := strings.Join(ss, ",")
	t.Logf("res1: %s", res1)

	is := []int{1, 2, 45}
	resInt := Join(",", is...)
	t.Logf("resInt: %s", resInt)

	i32s := []int32{1, 2, 45}
	resInt32 := Join(",", i32s...)
	t.Logf("resInt32: %s", resInt32)
}
```

### utils.Acquire基于原子操作的超时锁

```go
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
```
