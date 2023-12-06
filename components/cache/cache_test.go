package cache

import (
	"bytes"
	"fmt"
	"golang.org/x/exp/rand"
	"math"
	"runtime"
	"testing"
	"time"

	"github.com/grpc-boot/base/v2/kind"
	"github.com/grpc-boot/base/v2/kind/msg"
	"github.com/grpc-boot/base/v2/utils"
)

var (
	localDir = "/tmp/cache"
)

type User struct {
	Id        uint32 `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"createdAt"`
}

func (u User) ToMap() msg.Map {
	return msg.Map{
		"id":        u.Id,
		"name":      u.Name,
		"createdAt": u.CreatedAt,
	}
}

func TestCache_CommonMap(t *testing.T) {
	var (
		cache        = New(localDir, time.Second*3)
		id    uint32 = 10086
		key          = fmt.Sprintf("user:%d", 10086)
	)

	mp, err := CommonGet[msg.Map](cache, key, 60, func() (msg.Map, error) {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1008)))

		user := User{
			Id:        id,
			Name:      "移动",
			CreatedAt: time.Now().Unix(),
		}

		return user.ToMap(), nil
	})

	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	userMap := msg.MsgMap(mp)
	if userMap.Int("id") != int64(id) {
		t.Fatalf("want %d, got %v", id, userMap.Int("id"))
	}

	t.Logf("createdAt:%d", userMap.Int("createdAt"))
}

func TestCache_Get(t *testing.T) {
	start := time.Now()
	cache := New(localDir, time.Second*3)

	t.Logf("load from local cost:%s load length:%d", time.Since(start), Length(cache))

	var (
		effective bool
		data      = msg.Map{
			"bool":   true,
			"int":    100,
			"uint":   uint(math.MaxUint16),
			"uint64": uint64(math.MaxUint64),
			"float":  3.1415926,
			"string": `str测试`,
			"slice":  []interface{}{`string0`, `string1`, 1, 5},
			"bytes":  []byte(`bytes测试`),
			"time":   time.Now(),
		}
	)

	_ = Set(cache, "bool", data["bool"].(bool))
	_ = Set(cache, "int", data["int"].(int))
	_ = Set(cache, "uint", data["uint"].(uint))
	_ = Set(cache, "uint64", data["uint64"].(uint64))
	_ = Set(cache, "float", data["float"].(float64))
	_ = Set(cache, "string", data["string"].(string))
	_ = Set(cache, "slice", data["slice"].([]interface{}))
	_ = Set(cache, "bytes", data["bytes"].([]byte))
	_ = Set(cache, "time", data["time"].(time.Time))

	boolVal, effective, _ := Get(cache, "bool", 60)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	if msg.Bool(boolVal) != data.Bool("bool") {
		t.Fatalf("want %v, got %v", data.Bool("bool"), boolVal)
	}

	intVal, effective, _ := Get(cache, "int", 60)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	if msg.Int(intVal) != data.Int("int") {
		t.Fatalf("want %v, got %v", data.Int("int"), intVal)
	}

	uintVal, effective, _ := Get(cache, "uint", 60)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	if msg.Uint(uintVal) != data.Uint("uint") {
		t.Fatalf("want %v, got %v", data.Uint("uint"), uintVal)
	}

	uint64Val, effective, _ := Get(cache, "uint64", 60)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	if msg.Uint(uint64Val) != data.Uint("uint64") {
		t.Fatalf("want %v, got %v", data.Uint("uint64"), uint64Val)
	}

	floatVal, effective, _ := Get(cache, "float", 60)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	if msg.Float(floatVal) != data.Float("float") {
		t.Fatalf("want %v, got %v", data.Float("float"), floatVal)
	}

	stringValue, effective, _ := Get(cache, "string", 60)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	if msg.String(stringValue) != data.String("string") {
		t.Fatalf("want %v, got %v", data.String("string"), stringValue)
	}

	bytesValue, effective, _ := Get(cache, "bytes", 600)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	if !bytes.Equal(msg.Bytes(bytesValue), data.Bytes("bytes")) {
		t.Fatalf("want %s, got %v", data.Bytes("bytes"), bytesValue)
	}

	sliceValue, effective, _ := Get(cache, "slice", 60)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	t.Logf("want %v, got %v", data["slice"], sliceValue)

	timeValue, effective, _ := Get(cache, "time", 600)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	if !msg.Time(timeValue).Equal(data.Time("time")) {
		t.Fatalf("want %v, got %v", data.Time("time"), timeValue)
	}

	start = time.Now()
	cache.SyncLocal()
	t.Logf("sync cache to local cost: %s", time.Since(start))

	start = time.Now()
	cache = New(localDir, time.Second*3)

	t.Logf("load from local cost:%s load length:%d", time.Since(start), Length(cache))

	boolVal, effective, _ = Get(cache, "bool", 60)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	if msg.Bool(boolVal) != data.Bool("bool") {
		t.Fatalf("want %v, got %v", data.Bool("bool"), boolVal)
	}

	intVal, effective, _ = Get(cache, "int", 60)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	if msg.Int(intVal) != data.Int("int") {
		t.Fatalf("want %v, got %v", data.Int("int"), intVal)
	}

	uintVal, effective, _ = Get(cache, "uint", 60)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	if msg.Uint(uintVal) != data.Uint("uint") {
		t.Fatalf("want %v, got %v", data.Uint("uint"), uintVal)
	}

	uint64Val, effective, _ = Get(cache, "uint64", 60)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	if msg.Uint(uint64Val) != data.Uint("uint64") {
		t.Fatalf("want %v, got %v", data.Uint("uint64"), uint64Val)
	}

	floatVal, effective, _ = Get(cache, "float", 60)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	if msg.Float(floatVal) != data.Float("float") {
		t.Fatalf("want %v, got %v", data.Float("float"), floatVal)
	}

	stringValue, effective, _ = Get(cache, "string", 60)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	if msg.String(stringValue) != data.String("string") {
		t.Fatalf("want %v, got %v", data.String("string"), stringValue)
	}

	bytesValue, effective, _ = Get(cache, "bytes", 600)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	if !bytes.Equal(msg.Bytes(bytesValue), data.Bytes("bytes")) {
		t.Fatalf("want %s, got %v", data.Bytes("bytes"), bytesValue)
	}

	sliceValue, effective, _ = Get(cache, "slice", 60)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	t.Logf("want %v, got %v", data["slice"], sliceValue)

	timeValue, effective, _ = Get(cache, "time", 600)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	if !msg.Time(timeValue).Equal(data.Time("time")) {
		t.Fatalf("want %v, got %v", data.Time("time"), timeValue)
	}

	Delete(cache, "time")

	info, _ := utils.JsonEncode(InfoCache(cache))
	t.Logf("%s", info)
}

func TestCache_GetMap(t *testing.T) {
	start := time.Now()
	cache := New(localDir, time.Second*3)

	t.Logf("load from local cost:%s load length:%d", time.Since(start), Length(cache))

	var (
		key  = "map"
		data = msg.Map{
			"bool":    true,
			"int":     100,
			"uint64":  uint64(math.MaxUint64),
			"float":   3.1415926,
			"string":  `str测试`,
			"bytes":   []byte(`bytes测试`),
			"slice":   []interface{}{`string0`, `string1`, 1, 5},
			"ints":    []int{12, 9834234, 234234},
			"strings": []string{"12", "9834234", "234234"},
			"time":    time.Now(),
		}
	)

	err := Set(cache, key, data)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	cache.SyncLocal()

	cache = New(localDir, time.Second*3)

	value, effective, _ := Get(cache, key, 1)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}

	mp := msg.MsgMap(value)

	if mp.Bool("bool") != data.Bool("bool") {
		t.Fatalf("want %t, got %v", data.Bool("bool"), mp["bool"])
	}

	if mp.Int("int") != data.Int("int") {
		t.Fatalf("want %d, got %v", data.Int("int"), mp["int"])
	}

	if mp.Uint("uint64") != data.Uint("uint64") {
		t.Fatalf("want %d, got %v", data.Uint("uint64"), mp["uint64"])
	}

	if mp.Float("float") != data.Float("float") {
		t.Fatalf("want %v, got %v", data.Float("float"), mp["float"])
	}

	if mp.String("string") != data.String("string") {
		t.Fatalf("want %s, got %v", data.String("string"), mp["string"])
	}

	if !bytes.Equal(mp.Bytes("bytes"), data.Bytes("bytes")) {
		t.Fatalf("want %s, got %v", data.Bytes("bytes"), mp["bytes"])
	}

	t.Logf("want %v, got %v", data.Slice("slice"), mp.Slice("slice"))

	if !mp.Time("time").Equal(data.Time("time")) {
		t.Fatalf("want %v, got %v", data.Time("time"), mp.Time("time"))
	}

	if keyEqual(data.Ints("ints"), mp.Ints("ints")) {
		t.Logf("want %v, got %v", data.Ints("ints"), mp["ints"])
	}

	if keyEqual(data.Strings("strings"), mp.Strings("strings")) {
		t.Logf("want %v, got %v", data.Strings("strings"), mp["strings"])
	}
}

func keyEqual[T comparable](a, b kind.Slice[T]) bool {
	if a == nil && b == nil {
		return true
	}

	if len(a) != len(b) {
		return false
	}

	for index, _ := range a {
		if a[index] != b[index] {
			return false
		}
	}

	return true
}

func TestCache_SyncLocal(t *testing.T) {
	var (
		start = time.Now()
		cache = New(localDir, time.Second*3)
	)

	t.Logf("load %d length key cost:%s", Length(cache), time.Since(start))

	keyMax := 50 * 10000
	for i := 0; i < keyMax; i++ {
		key := fmt.Sprintf("key:%d", i)
		err := Set(cache, key, User{
			Id:        uint32(i),
			Name:      key,
			CreatedAt: time.Now().Unix(),
		}.ToMap())
		if err != nil {
			t.Fatalf("want nil, got %v", err)
		}
	}

	start = time.Now()
	cache.SyncLocal()

	t.Logf("sync to local cost: %s", time.Since(start))

	start = time.Now()
	runtime.GC()

	t.Logf("gc cost: %s", time.Since(start))
}

func TestCache_LoadLocal(t *testing.T) {
	var (
		start = time.Now()
		cache = New(localDir, time.Second*3)
	)

	t.Logf("load %d length key cost:%s", Length(cache), time.Since(start))

	start = time.Now()
	runtime.GC()

	t.Logf("first gc cost: %s", time.Since(start))

	start = time.Now()
	runtime.GC()

	t.Logf("second gc cost: %s", time.Since(start))
}

// BenchmarkCache_CommonMap-8       9286664               125.5 ns/op             0 B/op          0 allocs/op
func BenchmarkCache_CommonMap(b *testing.B) {
	var (
		cache        = New(localDir, time.Second*3)
		id    uint32 = 10086
		key          = fmt.Sprintf("user:%d", 10086)
	)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		value, err := CommonGet[msg.Map](cache, key, 10, func() (msg.Map, error) {
			time.Sleep(time.Millisecond * 500)

			user := User{
				Id:        id,
				Name:      "移动",
				CreatedAt: time.Now().Unix(),
			}

			return user.ToMap(), nil
		})

		if err != nil {
			b.Fatalf("want nil, got %v", err)
		}

		user := msg.MsgMap(value)

		if user.Int("id") != int64(id) {
			b.Fatalf("want %d, got %v", id, user.Int("id"))
		}
	}
}

// BenchmarkCacheParallel_CommonMap-8      20719171                58.49 ns/op            0 B/op          0 allocs/op
func BenchmarkCacheParallel_CommonMap(b *testing.B) {
	var (
		cache        = New(localDir, time.Second*3)
		id    uint32 = 10086
		key          = fmt.Sprintf("user:%d", 10086)
	)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			value, err := CommonGet[msg.Map](cache, key, 10, func() (msg.Map, error) {
				time.Sleep(time.Millisecond * 500)

				user := User{
					Id:        id,
					Name:      "移动",
					CreatedAt: time.Now().Unix(),
				}

				return user.ToMap(), nil
			})

			if err != nil {
				b.Fatalf("want nil, got %v", err)
			}

			user := msg.MsgMap(value)

			if user.Int("id") != int64(id) {
				b.Fatalf("want %d, got %v", id, user.Int("id"))
			}
		}
	})
}
