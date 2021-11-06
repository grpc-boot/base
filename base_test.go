package base

import (
	"hash/crc32"
	"testing"
	"time"
)

var (
	sm  ShardMap
	hs  HashSet
	btm Bitmap
)

type Data struct {
	CanHash

	id string
}

func (d *Data) HashCode() (hashValue uint32) {
	return crc32.ChecksumIEEE([]byte(d.id))
}

func init() {
	sm = NewShardMap()
	hs = NewHashSet(10)
	btm = NewBitmap(nil)
}

func TestHashValue(t *testing.T) {
	dd := &Data{
		id: "sfafd",
	}

	t.Fatal(dd.HashCode())
}

func TestMap(t *testing.T) {
	d := Data{
		id: "cc",
	}

	keyValue := map[interface{}]interface{}{
		"user": map[string]interface{}{
			"id":   15,
			"name": "ddadf",
		},
		"listLength": 34,
		"key":        "value",
		d:            55,
	}

	for key, value := range keyValue {
		sm.Set(key, value)
	}

	if int64(len(keyValue)) != sm.Length() {
		t.Fatalf("want %d, got %d", len(keyValue), sm.Length())
	}

	val, exists := sm.Get("user")
	if !exists {
		t.Fatalf("want true, got %t", exists)
	}

	if _, ok := val.(map[string]interface{}); !ok {
		t.Fatalf("want true, got %t", ok)
	}

	val, exists = sm.Get(d)
	if !exists {
		t.Fatalf("want true, got %t", exists)
	}

	if val != 55 {
		t.Fatal("want true, got false")
	}

	sm.Delete("key")

	if exists = sm.Exists("key"); exists {
		t.Fatalf("want false, got %t", exists)
	}

	if int64(len(keyValue)-1) != sm.Length() {
		t.Fatalf("want %d, got %d", len(keyValue)-1, sm.Length())
	}
}

// BenchmarkMap_SetParallel-4       3692262               330 ns/op              49 B/op          2 allocs/op
func BenchmarkMap_SetParallel(b *testing.B) {
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c := time.Now().UnixNano()
			sm.Set(c, c)
		}
	})
}

// BenchmarkMap_Set-4               2544157               641 ns/op             143 B/op          2 allocs/op
func BenchmarkMap_Set(b *testing.B) {
	b.ResetTimer()
	for index := 0; index < b.N; index++ {
		c := time.Now().UnixNano()
		sm.Set(c, c)
	}
}

// BenchmarkMap_GetParallel-4      14161170                95.0 ns/op             8 B/op          1 allocs/op
func BenchmarkMap_GetParallel(b *testing.B) {
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sm.Get(time.Now().UnixNano())
		}
	})
}

// BenchmarkMap_Get-4               5327881               225 ns/op               8 B/op          1 allocs/op
func BenchmarkMap_Get(b *testing.B) {
	b.ResetTimer()
	for index := 0; index < b.N; index++ {
		sm.Get(time.Now().UnixNano())
	}
}

func TestHashSet_Length(t *testing.T) {
	length := hs.Length()
	if length != 0 {
		t.Fatalf("length want 0, got %d", length)
	}

	newNum := hs.Add("1", 1000)
	if newNum != 2 {
		t.Fatalf("newNum want 2, got %d", newNum)
	}

	length = hs.Length()
	if length != 2 {
		t.Fatalf("length want 2, got %d", length)
	}

	newNum = hs.Add(1, 1000)
	if newNum != 1 {
		t.Fatalf("newNum want 1, got %d", newNum)
	}

	exists := hs.Exists(1000)
	if !exists {
		t.Fatal("want true, got false")
	}

	exists = hs.Exists("1000")
	if exists {
		t.Fatal("want false, got true")
	}
}

func TestBitmap_AddTag(t *testing.T) {
	t.Logf("hasBit:%t bitCount:%d binary:%s", btm.HasBit(), btm.BitCount(), btm.SprinfBinary())

	btm.AddTag(56)
	t.Logf("hasBit:%t bitCount:%d binary:%s", btm.HasBit(), btm.BitCount(), btm.SprinfBinary())

	btm.AddTag(8)
	t.Logf("hasBit:%t bitCount:%d binary:%s", btm.HasBit(), btm.BitCount(), btm.SprinfBinary())

	btm.AddTag(15)
	t.Logf("hasBit:%t bitCount:%d binary:%s", btm.HasBit(), btm.BitCount(), btm.SprinfBinary())

	btm.DelTag(15)
	t.Logf("hasBit:%t bitCount:%d binary:%s", btm.HasBit(), btm.BitCount(), btm.SprinfBinary())

	btm.DelTag(0)
	t.Logf("hasBit:%t bitCount:%d binary:%s", btm.HasBit(), btm.BitCount(), btm.SprinfBinary())
}
