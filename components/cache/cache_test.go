package cache

import (
	"bytes"
	"fmt"
	"golang.org/x/exp/rand"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/grpc-boot/base/v2/kind"
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

func (u User) ToMap() Map {
	return Map{
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

	user, err := cache.CommonMap(key, 60, func() (interface{}, error) {
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

	if user.Int("id") != int64(id) {
		t.Fatalf("want %d, got %v", id, user.Int("id"))
	}

	t.Logf("createdAt:%d", user.Int("createdAt"))
}

func TestCache_Get(t *testing.T) {
	start := time.Now()
	cache := New(localDir, time.Second*3)

	t.Logf("load from local cost:%s load length:%d", time.Since(start), cache.Length())

	var (
		effective bool
		data      = Map{
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

	for key, value := range data {
		err := cache.Set(key, value)
		if err != nil {
			t.Fatalf("want nil, got %v", err)
		}
	}

	boolVal, effective := cache.GetBool("bool", 60)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	if boolVal != data.Bool("bool") {
		t.Fatalf("want %v, got %v", data.Bool("bool"), boolVal)
	}

	intVal, effective := cache.GetInt("int", 60)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	if intVal != data.Int("int") {
		t.Fatalf("want %v, got %v", data.Int("int"), intVal)
	}

	uintVal, effective := cache.GetUint("uint", 60)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	if uintVal != data.Uint("uint") {
		t.Fatalf("want %v, got %v", data.Uint("uint"), uintVal)
	}

	uint64Val, effective := cache.GetUint("uint64", 60)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	if uint64Val != data.Uint("uint64") {
		t.Fatalf("want %v, got %v", data.Uint("uint64"), uint64Val)
	}

	floatVal, effective := cache.GetFloat("float", 60)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	if floatVal != data.Float("float") {
		t.Fatalf("want %v, got %v", data.Float("float"), floatVal)
	}

	stringValue, effective := cache.GetString("string", 60)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	if stringValue != data.String("string") {
		t.Fatalf("want %v, got %v", data.String("string"), stringValue)
	}

	bytesValue, effective := cache.GetBytes("bytes", 600)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	if !bytes.Equal(bytesValue, data.Bytes("bytes")) {
		t.Fatalf("want %s, got %v", data.Bytes("bytes"), bytesValue)
	}

	sliceValue, effective := cache.GetSlice("slice", 60)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	if !(reflect.DeepEqual(data.Slice("slice"), sliceValue)) {
		t.Fatalf("want %v, got %v", data["slice"], sliceValue)
	}

	timeValue, effective := cache.GetTime("time", 600)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	if !timeValue.Equal(data.Time("time")) {
		t.Fatalf("want %v, got %v", data.Time("time"), timeValue)
	}

	start = time.Now()
	cache.SyncLocal()
	t.Logf("sync cache to local cost: %s", time.Since(start))

	start = time.Now()
	cache = New(localDir, time.Second*3)

	t.Logf("load from local cost:%s load length:%d", time.Since(start), cache.Length())

	boolVal, effective = cache.GetBool("bool", 60)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	if boolVal != data.Bool("bool") {
		t.Fatalf("want %v, got %v", data.Bool("bool"), boolVal)
	}

	intVal, effective = cache.GetInt("int", 60)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	if intVal != data.Int("int") {
		t.Fatalf("want %v, got %v", data.Int("int"), intVal)
	}

	uintVal, effective = cache.GetUint("uint", 60)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	if uintVal != data.Uint("uint") {
		t.Fatalf("want %v, got %v", data.Uint("uint"), uintVal)
	}

	uint64Val, effective = cache.GetUint("uint64", 60)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	if uint64Val != data.Uint("uint64") {
		t.Fatalf("want %v, got %v", data.Uint("uint64"), uint64Val)
	}

	floatVal, effective = cache.GetFloat("float", 60)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	if floatVal != data.Float("float") {
		t.Fatalf("want %v, got %v", data.Float("float"), floatVal)
	}

	stringValue, effective = cache.GetString("string", 60)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	if stringValue != data.String("string") {
		t.Fatalf("want %v, got %v", data.String("string"), stringValue)
	}

	bytesValue, effective = cache.GetBytes("bytes", 600)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	if !bytes.Equal(bytesValue, data.Bytes("bytes")) {
		t.Fatalf("want %s, got %v", data.Bytes("bytes"), bytesValue)
	}

	sliceValue, effective = cache.GetSlice("slice", 60)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	t.Logf("want %v, got %v", data["slice"], sliceValue)

	timeValue, effective = cache.GetTime("time", 600)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}
	if !timeValue.Equal(data.Time("time")) {
		t.Fatalf("want %v, got %v", data.Time("time"), timeValue)
	}

	cache.Delete("time")

	info, _ := utils.JsonEncode(cache.Info())
	t.Logf("%s", info)
}

func TestCache_GetMap(t *testing.T) {
	start := time.Now()
	cache := New(localDir, time.Second*3)

	t.Logf("load from local cost:%s load length:%d", time.Since(start), cache.Length())

	var (
		key  = "map"
		data = Map{
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

	err := cache.Set(key, data)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	cache.SyncLocal()

	cache = New(localDir, time.Second*3)

	mp, effective := cache.GetMap(key, 1)
	if !effective {
		t.Fatalf("want true, got %v", effective)
	}

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
