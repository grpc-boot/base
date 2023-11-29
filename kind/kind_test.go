package kind

import (
	"math/rand"
	"strconv"
	"testing"
)

func TestSet_Add(t *testing.T) {
	hashSet := NewSet[int](23)
	if hashSet.Length() != 0 {
		t.Fatalf("want 0, got %d", hashSet.Length())
	}

	hashSet.Add(34)

	if !hashSet.Exists(34) {
		t.Fatalf("want true, got %v", hashSet.Exists(34))
	}

	if hashSet.Length() != 1 {
		t.Fatalf("want 1, got %d", hashSet.Length())
	}

	hashSet.Add(34)
	if hashSet.Length() != 1 {
		t.Fatalf("want 1, got %d", hashSet.Length())
	}

	hashSet.Add(56)
	if hashSet.Length() != 2 {
		t.Fatalf("want 2, got %d", hashSet.Length())
	}

	if !hashSet.Exists(56) {
		t.Fatalf("want true, got %v", hashSet.Exists(56))
	}

	hashSet.Del(56)
	if hashSet.Exists(56) {
		t.Fatalf("want false, got %v", hashSet.Exists(56))
	}
	if hashSet.Length() != 1 {
		t.Fatalf("want 1, got %d", hashSet.Length())
	}
}

func BenchmarkConcurrentSet_Add(b *testing.B) {

	hashSet := NewConcurrentSet[int32](16)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			val := rand.Int31n(2048)
			switch val % 3 {
			case 0:
				hashSet.Add(val)
			case 1:
				hashSet.Del(val)
			case 2:
				hashSet.Exists(val)
			}
		}
	})
}

func TestNewShardMap(t *testing.T) {
	var (
		sm  = NewShardMap[int64]()
		num = 102400
	)

	for index := 0; index < num; index++ {
		k := rand.Int63()
		sm.Set(k, true)
		if !sm.Exists(k) {
			t.Fatalf("want true, got %v", sm.Exists(k))
		}

		if rand.Uint32()%10 == 0 {
			delNum := sm.Delete(k)
			if delNum != 1 {
				t.Fatalf("want 1, got %v", delNum)
			}

			if sm.Exists(k) {
				t.Fatalf("want false, got %v", sm.Exists(k))
			}
		}
	}

	t.Logf("length:%d", sm.Length())
	t.Logf("shard length:%v", sm.ShardLength())
}

func TestNewConcurrentSet(t *testing.T) {
	d := -5
	c := uint32(d)

	t.Logf("%d", c)
}

func BenchmarkStringShardMap_Set(b *testing.B) {
	var (
		ssm = NewStringShardMap()
		key = `撒旦法`
	)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ssm.Set(key+strconv.Itoa(i), true)
	}
}

func BenchmarkShard_Set(b *testing.B) {
	var (
		ssm = NewShard[string](8)
		key = `撒旦法`
	)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ssm.Set(key+strconv.Itoa(i), true)
	}
}
