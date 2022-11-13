package optimistic

import (
	"testing"
	"time"
)

func BenchmarkAcquire(b *testing.B) {
	var lock int64

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			token := Acquire(&lock, time.Second)
			if token > 0 {
				time.Sleep(time.Second * 2)
				Release(&lock, token)
			}
		}
	})
}
